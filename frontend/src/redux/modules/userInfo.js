import produce from "immer"

// Actions
const LOGIN_SUCCESS = "userInfo/loginSuccess";


// Action creators
export function loginSuccess(email) {
  return {type: LOGIN_SUCCESS, payload: email};
}

// Initial state
const initialState = {authenticated: false, email: ""};

// Reducer
const reducer = produce((draft, action) => {
  switch(action.type) {
    case LOGIN_SUCCESS:
      draft.authenticated = true;
      draft.email = action.payload;
      break;
  }
}, initialState)

export default reducer;
