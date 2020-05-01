import ky from 'ky';

const api = ky.create({prefixUrl: "https://soesulcbkd.execute-api.us-east-1.amazonaws.com/alpha/", credentials: "include"})

export function getProject(projectId) {

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
