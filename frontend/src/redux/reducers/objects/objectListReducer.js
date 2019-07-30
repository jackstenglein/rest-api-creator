import { combineReducers } from 'redux';
import * as Actions from '../../actions/objects/objectListActions.js';
import {
    OBJECT_EDITOR_CANCEL,
    CREATE_OBJECT_SUCCESS
} from '../../actions/objects/objectEditorActions.js';

export default function objectList(state = {}, action) {
    switch (action.type) {
        case Actions.CLICK_CREATE:
            return {...state, create: true};
        case Actions.FETCH_OBJECTS_REQUEST:
            return {...state, status: Actions.FETCH_OBJECTS_REQUEST, errorMessage: null};
        case Actions.FETCH_OBJECTS_FAILURE:
            return {...state, status: Actions.FETCH_OBJECTS_FAILURE, errorMessage: action.message};
        case Actions.FETCH_OBJECTS_SUCCESS:
            return {...state, status: Actions.FETCH_OBJECTS_SUCCESS, errorMessage: null, items: action.objects};
        case Actions.CLICK_OBJECT:
            return {...state, selectedObject: action.id};
        case OBJECT_EDITOR_CANCEL:
        case CREATE_OBJECT_SUCCESS:
        case Actions.RESET_LIST_VIEW:
            return {...state, create: false, selectedObject: -1};
        default:
            return state;
    }
}
