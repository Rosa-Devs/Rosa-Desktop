import { GetAddres } from "../../wailsjs/go/app/App";
import { models } from "../models/manifest";

// Base URL of your API
let baseURL = ''; // Initialize baseURL as empty string

const initBaseURL = async () => {
  const address = await GetAddres(); // Get the address asynchronously
  baseURL = `http://${address}`; // Set the baseURL once the address is available
};

initBaseURL().then(() => {
  // Log the baseURL once it's initialized
  console.log(baseURL);
});
// let baseURL = 'http://localhost:8080'
console.log(baseURL)

// Function to handle HTTP requests
async function makeRequest(method: string, endpoint: string, data: any | null = null): Promise<any> {
  const requestOptions: RequestInit = {
    method: method,
    headers: {
      'Content-Type': 'application/json'
    }
  };

  if (data) {
    requestOptions.body = JSON.stringify(data);
  }

  try {
    const response = await fetch(baseURL + endpoint, requestOptions);
    if (!response.ok) {
      throw new Error('Request failed');
    }
    return await response.json();
  } catch (error) {
    console.error('Error making request:', error);
    throw error;
  }
}

// Authentication Endpoints
export interface createNewAccountProps {
  name: string
  avatar: string
}
export async function CreateNewAccount(props: createNewAccountProps): Promise<number | string> {
  return makeRequest('POST', '/aunth/create');
}

export async function Authorized(): Promise<boolean> {
  return makeRequest('GET', '/aunth/autorized');
}

export async function Trust(props: models.Message): Promise<number | string> {
  return makeRequest('POST', '/aunth/trust', props);
}

// Channels Endpoints
export interface createNewManifestProps {
  name: string,
  avatar: string,
}
export interface ManifestResponse {
  manifest: string;
}
export async function CreateNewManifest(props: createNewManifestProps): Promise<ManifestResponse | string> {
  return makeRequest('POST', '/channel/create', props);
}

export async function DeleteManifest(prosp: models.Manifest): Promise<number | string> {
  return makeRequest('POST', '/channel/delete', prosp);
}

export async function AddManifest(prosp: ManifestResponse): Promise<number | string> {
  return makeRequest('POST', '/channel/add', prosp);
}

export async function ListManifest(): Promise<models.Manifest[] | null> {
  return makeRequest('GET', '/channel/list');
}
// Events Endpoints
export async function ChangeListeningDb(props: models.Manifest): Promise<number | string> {
  return makeRequest('POST', '/event/change', props);
}

// Message Endpoints
export interface newMessageProps {
  msg: string
  manifest: models.Manifest
}
export async function NewMessage(props: newMessageProps): Promise<number | string> {
  return makeRequest('POST', '/message/new', props);
}

export async function MessagesList(props: models.Manifest): Promise<models.Message[] | null> {
  return makeRequest('POST', '/message/list', props);
}

