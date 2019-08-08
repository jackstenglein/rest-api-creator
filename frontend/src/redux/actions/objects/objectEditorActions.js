/**
* objectEditorActions.js
*
* Defines action types and action creators for the object editor.
*/

import fetch from 'cross-fetch';


/**
* Action types for creating/editing objects
*/

// Opening and using the object editor
export const OPEN_OBJECT_EDITOR_CREATE = "ObjectEditorCreate";
export const OBJECT_EDITOR_UPDATE_NAME = "ObjectEditorUpdateName";
export const OBJECT_EDITOR_UPDATE_DESCRIPTION = "ObjectEditorUpdateDescription";
export const OBJECT_EDITOR_ADD_ATTRIBUTE = "ObjectEditorAddAttribute";
export const OBJECT_EDITOR_REMOVE_ATTRIBUTE = "ObjectEditorRemoveAttribute";
export const OBJECT_EDITOR_UPDATE_ATTRIBUTE = "ObjectEditorUpdateAttribute";
export const OBJECT_EDITOR_INVALIDATE_ATTRIBUTE = "ObjectEditorInvalidateAttribute";
export const OBJECT_EDITOR_INVALIDATE_DETAILS = "ObjectEditorInvalidateDetails";
export const OBJECT_EDITOR_CLOSE_ERROR_MODAL = "ObjectEditorCloseErrorModal";
export const OBJECT_EDITOR_CANCEL = "ObjectEditorCancel";

// Object editor statuses
export const OBJECT_EDITOR_EDITING = "ObjectEditorEditing";
export const OBJECT_EDITOR_REQUEST_PENDING = "ObjectEditorRequestPending";
export const OBJECT_EDITOR_REQUEST_SUCCEEDED = "ObjectEditorRequestSucceeded";
export const OBJECT_EDITOR_REQUEST_FAILED = "ObjectEditorRequestFailed";

// Actual API requests for creating an object
export const CREATE_OBJECT_REQUEST = "CreateObjectRequest";
export const CREATE_OBJECT_FAILURE = "CreateObjectFailure";
export const CREATE_OBJECT_SUCCESS = "CreateObjectSuccess";

/**
* Other constants for creating/editing objects
*/

export const OBJECT_EDITOR_NOT_OPEN = "ObjectEditorNotOpen";
export const NEW_OBJECT_MODE = "ObjectEditorNew";


/**
* Action creators for using the object editor
*/

export function objectEditorUpdateName(name) {
    return {
        type: OBJECT_EDITOR_UPDATE_NAME,
        details: {
            name: name
        }
    };
}

export function objectEditorUpdateDescription(description) {
    return {
        type: OBJECT_EDITOR_UPDATE_DESCRIPTION,
        details: {
            description: description
        }
    };
}

export function objectEditorAddAttribute() {
    return {
        type: OBJECT_EDITOR_ADD_ATTRIBUTE
    };
}

export function objectEditorUpdateAttribute(index, update) {
    return {
        type: OBJECT_EDITOR_UPDATE_ATTRIBUTE,
        index: index,
        update: update
    };
}

export function objectEditorRemoveAttribute(index) {
    return {
        type: OBJECT_EDITOR_REMOVE_ATTRIBUTE,
        index: index
    };
}

export function objectEditorInvalidateAttribute(index, feedback) {
    return {
        type: OBJECT_EDITOR_INVALIDATE_ATTRIBUTE,
        index: index,
        feedback: feedback
    }
}

export function objectEditorInvalidateDetails(feedback) {
    return {
        type: OBJECT_EDITOR_INVALIDATE_DETAILS,
        details: {
            ...feedback
        }
    }
}

export function objectEditorCloseErrorModal() {
    return {type: OBJECT_EDITOR_CLOSE_ERROR_MODAL};
}

export function objectEditorCancel() {
    return {type: OBJECT_EDITOR_CANCEL};
}

/**
* Action creators for creating objects
*/

function requestCreateObject(object) {
    return {
        type: CREATE_OBJECT_REQUEST,
        object
    }
}

function createObjectResponse(object, json) {
    console.log("JsonError: ", json.error);
    if (json.error !== undefined) {
        return {
            type: CREATE_OBJECT_FAILURE,
            error: json.error,
            object
        }
    }

    console.log("Emitting success action");
    return {
        type: CREATE_OBJECT_SUCCESS,
        object: object
    }
}

function createObjectFailure(message, feedback) {
    return {
        type: CREATE_OBJECT_FAILURE,
        message: message,
        feedback
    }
}

export function createObject(object) {
    console.log("OBJECT: ", object);
    return function(dispatch) {
        dispatch(requestCreateObject(object));
        return fetch('http://localhost:8000/app/projects/5/objects', {
            method: "POST",
            body: JSON.stringify(object),
            headers: {
                "Content-Type": "application/json"
            },
            credentials: "include"
        }).then(function(response, error) {
            console.log("Response: ", response);
            if (error) {
                console.log('An error ocurred.', error);
            }
            return response.json();
        }).then(function(json) {
            console.log("Json response: ", json);
            dispatch(createObjectResponse(object, json));
        }).catch(function(error) {
            console.log("Caught error: ", error);
            let message = "Unable to create object \"" + object.name + "\" due to an" +
                " unknown error. Please check your internet connection and try again.";
            dispatch(createObjectFailure(message, null));
        });
    }
}