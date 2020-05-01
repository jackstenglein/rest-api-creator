import ky from 'ky';

const api = ky.create({prefixUrl: "https://soesulcbkd.execute-api.us-east-1.amazonaws.com/alpha/"})

export function getProject(projectId) {

}

export async function login(email, password) {
  try {
    const response = await api.put("login", {json: {email: email, password: password}, credentials: "include"}).json();
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
    const response = await api.put(`projects/${projectId}/objects`, {json: object, credentials: "include"}).json();
    return response;
  } catch (err) {
    if (err.response !== undefined) {
      console.log(await err.response.json());
    }
    return undefined;
  }
}

export function signup(username, password) {
  
}
