import { combineReducers } from 'redux';
import * as Actions from '../../actions/objects/objectEditorActions.js';

// Reducer that handles the name and description section of the object editor
// State is details section of editor object
function objectEditorDetails(state = {}, action) {
    switch (action.type) {
        case Actions.OBJECT_EDITOR_UPDATE_NAME:
        case Actions.OBJECT_EDITOR_UPDATE_DESCRIPTION:
        case Actions.OBJECT_EDITOR_INVALIDATE_DETAILS:
            return Object.assign({}, state, action.details);
        default:
            return state;
    }
}

// Reducer that handles attributes section of the object editor
// State is list of attributes
function objectEditorAttributes(state = [], action) {
    switch (action.type) {
        case Actions.OBJECT_EDITOR_ADD_ATTRIBUTE:
            return  [...state, {
                name: "",
                nameFeedback: null,
                type: "Choose...",
                typeFeedback: null,
                default: "",
                defaultFeedback: null
            }];
        case Actions.OBJECT_EDITOR_UPDATE_ATTRIBUTE:
            return state.map((attribute, index) => {
                if (index === action.index) {
                    return Object.assign({}, attribute, action.update);
                }
                return attribute;
            });
        case Actions.OBJECT_EDITOR_INVALIDATE_ATTRIBUTE:
            return state.map((attribute, index) => {
                if (index === action.index) {
                    return Object.assign({}, attribute, action.feedback);
                }
                return attribute;
            });
        case Actions.OBJECT_EDITOR_REMOVE_ATTRIBUTE:
            return [...state.slice(0, action.index), ...state.slice(action.index + 1)];
        default:
            return state;
    }
}

// Reducer that handles control section of the object editor
// State is control section of the editor object
function objectEditorControl(state = {}, action) {
    switch (action.type) {
        case Actions.CREATE_OBJECT_REQUEST:
            return {...state, status: Actions.OBJECT_EDITOR_REQUEST_PENDING};
        case Actions.CREATE_OBJECT_SUCCESS:
            return {...state, status: Actions.OBJECT_EDITOR_REQUEST_SUCCEEDED};
        case Actions.CREATE_OBJECT_FAILURE:
            return {...state, status: Actions.OBJECT_EDITOR_REQUEST_FAILED, errorMessage: action.message};
        case Actions.OBJECT_EDITOR_CLOSE_ERROR_MODAL:
            return {...state, status: Actions.OBJECT_EDITOR_EDITING, errorMessage: ""};
        default:
            return state;
    }
}

const objectEditor = combineReducers({
    details: objectEditorDetails,
    attributes: objectEditorAttributes,
    control: objectEditorControl
});

export default objectEditor;
