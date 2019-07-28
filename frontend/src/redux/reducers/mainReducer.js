import objectEditor from './objects/objectEditorReducer.js';

const initialState = {
    selectedProject: 5,
    objects: {
        editor: {
            control: {
                status: null,
                selectedObject: -1,
                errorMessage: ""
            },
            details: {
                name: "",
                description: "",
                nameFeedback: null,
            },
            attributes: []
        },
        items: []
    }
}



export function crudCreatorApp(state = initialState, action) {
    return {
        objects: {
            editor: objectEditor(state.objects.editor, action),
            items: []
        }
    }
}
