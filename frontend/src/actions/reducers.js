import * as Actions from './actions.js';

const initialState = {
    selectedProject: 5,
    objects: {
        isCreating: false,
        editor: {
            mode: Actions.OBJECT_EDITOR_NOT_OPEN,
            selectedObject: -1,
            name: "",
            attributes: []
        },
        items: []
    }
}



export function crudCreatorApp(state = initialState, action) {
    return {
        objects: objects(state.objects, action)
    }
}

function objects(state = {}, action) {
    let editor;

    switch (action.type) {
        case Actions.CREATE_OBJECT_REQUEST:
            return Object.assign({}, state, {
                isCreating: true
            });
        case Actions.CREATE_OBJECT_SUCCESS:
            return Object.assign({}, state, {
                isCreating: false,
                items: [
                    ...state.items,
                    action.object
                ]
            });
        case Actions.OBJECT_EDITOR_UPDATE_NAME:
            editor = Object.assign({}, state.editor, {
                name: action.name
            });
            return Object.assign({}, state, {
                editor: editor
            });
        case Actions.OBJECT_EDITOR_UPDATE_DESCRIPTION:
            editor = Object.assign({}, state.editor, {
                description: action.description
            });
            return Object.assign({}, state, {
                editor: editor
            });
        case Actions.OBJECT_EDITOR_ADD_ATTRIBUTE:
            editor = Object.assign({}, state.editor, {
                attributes: [
                    ...state.editor.attributes,
                    {
                        name: "",
                        type: "Choose...",
                        default: ""
                    }
                ]
            });
            return Object.assign({}, state, {
                editor: editor
            });
        case Actions.OBJECT_EDITOR_UPDATE_ATTRIBUTE:
            editor = Object.assign({}, state.editor, {
                attributes: state.editor.attributes.map((attribute, index) => {
                    if (index === action.index) {
                        return Object.assign({}, attribute, action.update);
                    } else {
                        return attribute;
                    }
                })
            });
            return Object.assign({}, state, {
                editor: editor
            });
        case Actions.OBJECT_EDITOR_REMOVE_ATTRIBUTE:
            let firstList = state.editor.attributes.slice(0, action.index);
            let secondList = state.editor.attributes.slice(action.index + 1);
            editor = Object.assign({}, state.editor, {
                attributes: [
                    ...firstList,
                    ...secondList
                ]
            });
            return Object.assign({}, state, {
                editor: editor
            });
        default:
            return state;
    }
}
