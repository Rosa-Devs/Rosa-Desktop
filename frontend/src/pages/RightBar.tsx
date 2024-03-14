import React, { useState, useEffect, ChangeEvent, createRef} from 'react';
import Buble from '../itemkit/ChatBuble'; // Import the Buble component if not already imported
import { CgAdd } from "react-icons/cg";
import { MdCreate } from "react-icons/md";
import { Cropper, ReactCropperElement } from 'react-cropper';
import 'cropperjs/dist/cropper.css';
import { toast } from 'react-toastify';
import Optional from '../models/Optional';
import { AddManifest, CreateNewManifest, ListManifest } from '../api/api';
import { models } from '../models/manifest';


const RightSidebar = ({ setManifest }: { setManifest: React.Dispatch<React.SetStateAction<any>> }) => {
  const [contacts, setContacts] = useState<models.Manifest[] | null>(null);
  const [isUploadPopupOpen, setIsUploadPopupOpen] = useState(false);
  const [isCreatePopupOpen, setIsCreatePopupOpen] = useState(false);

  useEffect(() => {
    const fetchContacts = async () => {
      try {
        const contactsData = await ListManifest();
        setContacts(contactsData);
      } catch (error) {
        console.error('Error fetching contacts:', error);
      }
      setTimeout(() => {
        fetchContacts();
      }, 2000);
    };

    fetchContacts();
  }, []);

  // useEffect(() => {
  //   if (contacts !== null) {
  //     setManifest(contacts[0])
  //   }
  // })

  const handleAddManifest = () => {
    setIsUploadPopupOpen(true);
  };

  const handleCreateManifest = () => {
    setIsCreatePopupOpen(true);
  }
  const handleCloseCreateManifest = () => {
    setIsCreatePopupOpen(true);
  }

  const handleCloseUploadPopup = () => {
    setIsUploadPopupOpen(false);
  };


  const handleFileUpload = async (file: File) => {
    // Handle file upload logic here
    // You can use this function to upload the selected file
    try {
      
      const fileContents = await readFile(file);
      console.log('File Data:', fileContents);

      // @ts-ignore
      AddManifets(fileContents)
      window.location.reload();
    } catch (error) {
      console.error('Error reading file:', error);
    }
    // After handling the file, close the popup
    handleCloseUploadPopup();
  };

  const readFile = (file: File) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = (event: any) => {
        resolve(event.target.result);
      };

      reader.onerror = (error) => {
        reject(error);
      };

      reader.readAsText(file);
    });
  };

  // const [nodes, setnodes] = useState(0);
  // useEffect(() => {
  //   const fetchNodes = async () => {
  //     try {
  //       const contactsData = await Nodes();
  //       setnodes(contactsData);
  //     } catch (error) {
  //       console.error('Error fetching contacts:', error);
  //     }
  //     setTimeout(() => {
  //       fetchNodes();
  //     }, 2000);
  //   };

  //   fetchNodes();
  // }, []);
  
  
  


  return (
    <div className="sidebar">
      {contacts !== null && contacts.length === 0 ? (
        <p>No rooms!</p>
      ) : (
        <div className="content">
          <div className="bubles">
            {contacts?.map((contact) => (
              <Buble key={contact.pubsub} contact={contact} setManifest={setManifest} />
            ))}
          </div>
          <div className="add-button-container">
            <button onClick={handleAddManifest} className="add-button"><CgAdd className='text-4xl'/></button><br />
            <button onClick={handleCreateManifest} className="add-button"><MdCreate className='text-4xl'/></button>
            {/* <p className='dht'>DHT: {nodes}</p> */}
          </div>
        </div>
      )}

      {isUploadPopupOpen && (
        <div className="upload-popup">
          <div className="upload-popup-content">
            <input type="file" className='upload-form' id='file-form' onChange={(e) => handleFileUpload(e.target.files![0])}/>
            <button onClick={handleCloseUploadPopup}>Close</button>
          </div>
        </div>
      )}

      {isCreatePopupOpen && (
        <CreateManifestPopUp close={setIsCreatePopupOpen}/>
      )}
    </div>
  );
};

export default RightSidebar;



interface CreateManifestPopUpProps {
  close: (value: boolean) => void;
}

const CreateManifestPopUp: React.FC<CreateManifestPopUpProps> = ({ close }) => {
  // ref of the file input
  const fileRef = createRef<HTMLInputElement>();

  // the selected image
  const [uploaded, setUploaded] = useState(null as string | null);

  // the resulting cropped image
  const [cropped, setCropped] = useState();

  // the reference of cropper element
  const cropperRef = createRef<ReactCropperElement>();

  const [name, setName] = useState("")


  const onFileInputChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
    const file = e.target?.files?.[0];
    if (file) {
      file2Base64(file).then((base64) => {
        setUploaded(base64);
      });
    }
  }

  const onCrop = () => {
    const imageElement: any = cropperRef?.current;
    const cropper: any = imageElement?.cropper;
    setCropped(cropper.getCroppedCanvas().toDataURL())
    return cropper.getCroppedCanvas().toDataURL()
  }


  const hanleCreateBtn = async () => {
    try {
      // Create manifest

      const jsonData: Optional = {
        Image: onCrop()
      };

 

      // Convert the JSON object to a JSON string
      const jsonString = JSON.stringify(jsonData);

      const manifest = await CreateNewManifest({name: name,avatar: jsonString})

      if (typeof manifest === 'string') {
        toast("Err")
      } else {
        AddManifest(manifest)
      }
      

      toast("New chanell created!")
      close(false)
    } catch (error) {
      // Handle the error and show a toast message
      toast("Avatar not specified");
    }
  }

  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFilePathChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      setSelectedFile(file);
    }
  };

  return (
    <div className="create-popup">
      <div className="create-icon-button">
        <div className="">
          <div className="">
            <div className="App">
              { uploaded ?
                  <div>
                    <Cropper
                      src={uploaded}
                      style={{height: 400, width: 400}}
                      autoCropArea={1}
                      aspectRatio={1}
                      viewMode={1}
                      guides={true}
                      ref={cropperRef}
                      className='mb-3'
                    />
                    {/* <button onClick={onCrop}>Crop</button> */}
                    {/* {cropped && <img src={cropped} alt="Cropped!"/>} */}
                  </div>
                  :
                  <>
                    <input
                      type="file"
                      style={{display: 'none'}}
                      ref={fileRef}
                      onChange={onFileInputChange}
                      accept="image/png,image/jpeg,image/gif"
                    />
                    <button
                      className='cropper-btn'
                      onClick={() => fileRef.current?.click()}
                    >Upload Avatar
                    </button>
                  </>}
            </div>
            <div className="">
              <input type="text" placeholder="Name" className="height"
                value={name}
                onChange={(event) => {setName(event.target.value)}}
              />
            </div>
            <div className="">
              
            </div>
          </div>
        </div>
      </div>
      <div className="">
        <button className='close-btn' onClick={() => close(false)}>Close</button>
        <button className='close-btn ml-2' onClick={hanleCreateBtn}>Create</button>
      </div>
    </div>
  );
};





// this transforms file to base64
const file2Base64 = (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result?.toString() || '');
    reader.onerror = (error) => reject(error);
  });
};

