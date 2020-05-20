// projects.js contains the actions, action creators and reducer for projects.

import produce from "immer"
import * as network from './network.js';
import { GET_USER_RESPONSE } from './userInfo.js';

// Actions
const FETCH_PROJECT_REQUEST = "projects/fetchProjectRequest";
const FETCH_PROJECT_RESPONSE = "projects/fetchProjectResponse";
const PUT_OBJECT_SUCCESS = "objects/putObjectSuccess";

// Action creators
export function fetchProjectRequest(projectId) {
  return {type: FETCH_PROJECT_REQUEST, payload: projectId};
}

export function fetchProjectResponse(projectId, response) {
  return {type: FETCH_PROJECT_RESPONSE, payload: {id: projectId, response: response}};
}

export function putObjectSuccess(projectId, object) {
  return {
    type: PUT_OBJECT_SUCCESS, 
    payload: {
      projectId: projectId, 
      originalId: object.id,
      object: {...object, id: object.name.toLowerCase()}
    }
  };
}

// Initial state
const initialState = {}

// Full projects reducer -- draft = projects = {pid1: {...}, pid2: {...}, ...}
const reducer = produce((draft, action = {}) => {
  switch (action.type) {
    case FETCH_PROJECT_REQUEST:
      draft[action.payload] = {network: network.pending()};
      break;
    case FETCH_PROJECT_RESPONSE:
      const response = action.payload.response;
      const id = action.payload.id
      if (response.error !== undefined) {
        draft[id] = {network: network.failure(response.error), objects: {}};
      } else {
        draft[id] = response.project;
        draft[id].network = network.success();
      }
      break;
    case PUT_OBJECT_SUCCESS:
      const object = action.payload.object;
      const originalId = action.payload.originalId;
      const projectId = action.payload.projectId;
      if (!draft[projectId].objects) {
        draft[projectId].objects = {};
      }
      if (originalId !== undefined && originalId !== object.id) {
        delete draft[projectId].objects[originalId];
      }
      draft[projectId].objects[object.id] = object; 
      break;
    case GET_USER_RESPONSE:
      if (action.payload.user && action.payload.user.projects) {
        Object.entries(action.payload.user.projects).forEach(([id, project]) => {
          draft[id] = project;
          draft[id].network = network.success();
        })
      }
      break;
  }
}, initialState)

export default reducer;
