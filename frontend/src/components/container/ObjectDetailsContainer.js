import { connect } from 'react-redux'
import ObjectDetails from '../presentation/ObjectDetails';
import {
    clickEdit,
    fetchObjectIfNeeded
} from '../../redux/actions/objects/objectViewActions.js';
import { OBJECT_EDITOR_EDITING } from '../../redux/actions/objects/objectEditorActions.js';


function getObject(objects, id) {
    for (let i = 0; i < objects.length; ++i) {
        if (objects[i].id === id) {
            return objects[i];
        }
    }
    return null;
}

function mapStateToProps(state, ownProps) {
    console.log("Own props: ", ownProps);
    console.log("state:", state);
    // return Object.assign({}, state.objects.list.items[31], {projectName: 'Test Project'});
    const id = parseInt(ownProps.match.params.id);
    return {
        selectedObject: id,
        object: getObject(state.objects.list.items, id),
        projectName: 'Test Project',
        redirectToEdit: state.objects.editor.control.status === OBJECT_EDITOR_EDITING
    };
}

const mapDispatchToProps = dispatch => {
    return {
        callbacks: {
            clickEdit: (object) => dispatch(clickEdit(object)),
            fetchObjectIfNeeded: (id) => dispatch(fetchObjectIfNeeded(id))
        }
    };
}

const ObjectDetailsContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectDetails);

export default ObjectDetailsContainer;
