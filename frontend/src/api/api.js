import ky from 'ky';

const config = {
  dev: {
    url: "https://soesulcbkd.execute-api.us-east-1.amazonaws.com/dev/"
  },
  alpha: {
    url: "https://9h52owapyf.execute-api.us-east-1.amazonaws.com/alpha/"
  }
};

const api = ky.create({prefixUrl: config[process.env.REACT_APP_STAGE].url, credentials: "include"});
const defaultError = {error: "Failed to make network request."};

export async function deleteObject(projectId, objectId) {
  try {
    return await api.delete(`projects/${projectId}/objects/${objectId}`).json();
  } catch (err) {
    if (err.response !== undefined) {
      return await err.response.json();
    }
    return defaultError;
  }
}

export async function getDownloadLink(projectId) {
  const timeout = 30000; // 30 seconds
  try {
    return await api.get(`projects/${projectId}/code`, {timeout: timeout}).json();
  } catch (err) {
    if (err.response !== undefined) {
      return await err.response.json();
    }
    return defaultError;
  }
}

export async function getProject(projectId) {
  try  {
    return await api.get(`projects/${projectId}`).json();
  } catch (err) {
    if (err.response !== undefined) {
      return await err.response.json();
    }
    return defaultError;
  }
}

export async function getUser() {
  try {
    return await api.get(`user`).json();
  } catch (err) {
    if (err.response !== undefined) {
      return await err.response.json();
    }
    return defaultError;
  }
}

export async function login(email, password) {
  try {
    const response = await api.put("login", {json: {email: email, password: password}}).json();
    return response;
  } catch (err) {
    console.log("Caught error: ", err);
    if (err.response !== undefined) {
      console.log("Getting response json");
      return await err.response.json();
    }
    return undefined;
  }
}

export function logout() {

}

export async function putObject(projectId, object) {
  try {
    const response = await api.put(`projects/${projectId}/objects`, {json: object}).json();
    return response;
  } catch (err) {
    if (err.response !== undefined) {
      console.log(await err.response.json());
    }
    return undefined;
  }
}

export async function signup(email, password) {
  try {
    return await api.post("signup", {json: {email: email, password: password}}).json()
  } catch (err) {
    if (err.response !== undefined) {
      return await err.response.json();
    }
    return undefined;
  }
}
