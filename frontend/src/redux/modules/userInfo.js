import produce from "immer"
import * as network from './network.js';

// Actions
const LOGIN_SUCCESS = "userInfo/loginSuccess";
const GET_USER_REQUEST = "userInfo/getUserRequest";
export const GET_USER_RESPONSE = "userInfo/getUserResponse";


// Action creators
export function loginSuccess(email) {
  return {type: LOGIN_SUCCESS, payload: email};
}

export function getUserRequest() {
  return {type: GET_USER_REQUEST};
}

export function getUserResponse(response) {
  return {type: GET_USER_RESPONSE, payload: response};
}

// Initial state
const initialState = {
  authenticated: false, 
  email: "",
  network: network.none()
};

// Authenticated reducer
function authReducer(state, action) {
  switch (action.type) {
    case LOGIN_SUCCESS:
      return true;
    case GET_USER_REQUEST:
      return state;
    case GET_USER_RESPONSE:
      return action.payload.error === undefined;
    default:
      return state;
  }
}

// Email reducer
function emailReducer(state, action) {
  switch (action.type) {
    case LOGIN_SUCCESS:
      return action.payload;
    case GET_USER_RESPONSE:
      if (action.payload.error) {
        return "";
      }
      return action.payload.email;
    default:
      return state;
  }
}

// Network reducer
function networkReducer(state, action) {
  switch (action.type) {
    case GET_USER_REQUEST:
      return network.pending();
    case GET_USER_RESPONSE:
      if (action.payload.error) {
        return network.failure(action.payload.error);
      }
      return network.success();
    default:
      return state;
  }
}

// Main reducer
const reducer = produce((draft, action) => {
  draft.authenticated = authReducer(draft.authenticated, action);
  draft.email = emailReducer(draft.email, action);
  draft.network = networkReducer(draft.network, action);
}, initialState)

export default reducer;
