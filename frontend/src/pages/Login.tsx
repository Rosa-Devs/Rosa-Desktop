import { ChangeEvent, createRef, useState } from "react";
import { Cropper, ReactCropperElement } from 'react-cropper';
import { toast } from "react-toastify";
import { CreateNewAccount } from "../../wailsjs/go/core/Core";



const Login = () => {
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
      if (name === "") {
        toast("You need to specify a nickname")
        return
      }

      await CreateNewAccount(name, onCrop())

      toast("Hello im jarvis.... Just kidding..")
      toast("Please wait a sec we are doing some cryptography :)")

      setTimeout(()=>{
        window.location.reload()
      }, 500)

      
  }

  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFilePathChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      setSelectedFile(file);
    }
  };

  return (
    <div className="create-popup calm-gradient-background">
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
              <input type="text" placeholder="Nickname" className="height"
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
        <button className='close-btn ml-2' onClick={hanleCreateBtn}>Create</button>
      </div>
    </div>
  );
}

const file2Base64 = (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result?.toString() || '');
    reader.onerror = (error) => reject(error);
  });
};

export default Login