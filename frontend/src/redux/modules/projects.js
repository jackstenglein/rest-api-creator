// projects.js contains the actions, action creators and reducer for projects.

import produce from "immer"

// Actions
const PUT_PROJECT_SUCCESS = "projects/putProjectSuccess";
const PUT_OBJECT_SUCCESS = "objects/putObjectSuccess";


// Action creators
export function putProjectSuccess(project) {
  return {type: PUT_PROJECT_SUCCESS, payload: project};
}

export function putObjectSuccess(projectId, object) {
  return {type: PUT_OBJECT_SUCCESS, payload: {projectId: projectId, object: object}};
}

// Initial state
const initialState = {}

// Full projects reducer -- draft = projects = {pid1: {...}, pid2: {...}, ...}
const reducer = produce((draft, action = {}) => {
  switch (action.type) {
    case PUT_PROJECT_SUCCESS:
      draft[action.payload.id] = action.payload;
      break;
    case PUT_OBJECT_SUCCESS:
      const object = action.payload.object;
      const projectId = action.payload.projectId
      draft[projectId].objects[object.id] = object; 
      break;
  }
}, initialState)

export default reducer;
