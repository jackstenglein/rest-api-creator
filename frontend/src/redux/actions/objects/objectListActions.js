/**
* objectListActions.js
*
* Defines action types and action creators for the object list view.
*/

import fetch from 'cross-fetch';


/**
* Action types
*/

export const CLICK_CREATE = "ObjectListViewClickCreate";
export const CLICK_OBJECT = "ObjectListViewClickObject";
export const RESET_LIST_VIEW = "ObjectListViewReset";

export const FETCH_OBJECTS_REQUEST = "FetchObjectsRequest";
export const FETCH_OBJECTS_FAILURE = "FetchObjectsFailure";
export const FETCH_OBJECTS_SUCCESS = "FetchObjectsSuccess";


/**
* Action creators
*/

export function clickCreate() {
    return {type: CLICK_CREATE};
}

export function clickObject(id) {
    return {
        type: CLICK_OBJECT,
        id: id
    };
}

export function resetListView() {
    return {type: RESET_LIST_VIEW};
}

function requestFetchObjects() {
    return {
        type: FETCH_OBJECTS_REQUEST,
    };
}

function fetchObjectsResponse(json) {
    if (json.error !== undefined) {
        console.log("JsonError: ", json.error);
        return {
            type: FETCH_OBJECTS_FAILURE,
            error: json.error
        }
    }

    console.log("Emitting success action");
    return {
        type: FETCH_OBJECTS_SUCCESS,
        objects: json.objects
    };
}

function fetchObjectsFailure(message) {
    return {
        type: FETCH_OBJECTS_FAILURE,
        message: message
    };
}

function shouldFetchObjects(state) {
    return state.objects.list.status !== FETCH_OBJECTS_REQUEST &&
           state.objects.list.status !== FETCH_OBJECTS_SUCCESS;
}

export function fetchObjects(project=5) {
    return function(dispatch) {
        dispatch(requestFetchObjects());
        return fetch('http://localhost:8000/app/projects/' + project + '/objects', {
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
            dispatch(fetchObjectsResponse(json));
        }).catch(function(error) {
            console.log("Caught error: ", error);
            let message = "Unable to find your objects due to an unknown error. " +
                "Please check your internet connection and try again.";
            dispatch(fetchObjectsFailure(message));
        });
    }
}

export function fetchObjectsIfNeeded() {
    return function(dispatch, getState) {
        const state = getState();
        if (shouldFetchObjects(state)) {
            console.log("Fetching objects");
            return dispatch(fetchObjects(state.selectedProject));
        } else {
            console.log("Already fetching");
            return Promise.resolve();
        }
    }
}
