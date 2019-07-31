/**
* objectViewActions.js
*
* Defines action types and action creators for the single object view.
*/

import fetch from 'cross-fetch';


/**
* Action types
*/

export const CLICK_EDIT = "ObjectViewClickEdit";
export const CLICK_DELETE = "ObjectViewClickDelete";

export const FETCH_OBJECT_REQUEST = "ObjectViewFetchObjectRequest";
export const FETCH_OBJECT_FAILURE = "ObjectViewFetchObjectFailure";
export const FETCH_OBJECT_SUCCESS = "ObjectViewFetchObjectSuccess";


/**
* Action creators
*/

export function clickEdit(object) {
    return {
        type: CLICK_EDIT,
        object: object
    };
}

// function requestFetchObject() {
//     return {
//         type: FETCH_OBJECT_REQUEST
//     };
// }

function fetchObject(project=5, object) {
    return function(dispatch) {
        // dispatch(requestFetchObjects());
        return fetch('http://localhost:8000/app/projects/' + project + '/objects/' + object, {
            method: "GET",
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
            // dispatch(fetchObjectsResponse(json));
        }).catch(function(error) {
            console.log("Caught error: ", error);
            let message = "Unable to find your objects due to an unknown error. " +
                "Please check your internet connection and try again.";
            // dispatch(fetchObjectsFailure(message));
        });
    }
}

function getObject(objects, id) {
    for (let i = 0; i < objects.length; ++i) {
        if (objects[i].id === id) {
            return objects[i];
        }
    }
    return null;
}

function shouldFetchObject(state, object) {
    return getObject(state.objects.list.items, object) === null;
}


export function fetchObjectIfNeeded(object) {
    return function(dispatch, getState) {
        const state = getState();
        if (shouldFetchObject(state, object)) {
            console.log("Fetching object %d", object);
            return dispatch(fetchObject(state.selectedProject, object));
        } else {
            console.log("Already fetching or cached");
            return Promise.resolve();
        }
    }
}
