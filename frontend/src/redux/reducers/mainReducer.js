import objectEditor from './objects/objectEditorReducer.js';
import objectList from './objects/objectListReducer.js';

const initialState = {
    selectedProject: 5,
    objects: {
        list: {
            status: null,
            errorMessage: null,
            create: false,
            selectedObject: -1,
            items: []
        }
    }
}



export function crudCreatorApp(state = initialState, action) {
    return {
        objects: {
            editor: objectEditor(state.objects.editor, action),
            list: objectList(state.objects.list, action)
        }
    }
}
