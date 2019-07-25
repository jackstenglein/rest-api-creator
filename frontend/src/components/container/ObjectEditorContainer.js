import { connect } from 'react-redux';
import ObjectEditor  from '../presentation/ObjectEditor';
import {
    objectEditorUpdateName,
    objectEditorUpdateDescription,
    objectEditorAddAttribute,
    objectEditorUpdateAttribute,
    objectEditorRemoveAttribute
}  from '../../actions/actions';

const testProps = {
    projectName: 'Test Project'
}

const mapStateToProps = (state, ownProps) => {
    return Object.assign({}, state.objects.editor, testProps);
}

function nameOnChange(dispatch, event) {
    dispatch(objectEditorUpdateName(event.target.value));
}

function descriptionOnChange(dispatch, event) {
    dispatch(objectEditorUpdateDescription(event.target.value));
}

function clickAddAttribute(dispatch) {
    dispatch(objectEditorAddAttribute());
}

function updateAttribute(dispatch, index, update) {
    dispatch(objectEditorUpdateAttribute(index, update));
}

function removeAttribute(dispatch, index) {
    dispatch(objectEditorRemoveAttribute(index));
}

const mapDispatchToProps = dispatch => {
    return {
        nameOnChange: event => nameOnChange(dispatch, event),
        descriptionOnChange: event => descriptionOnChange(dispatch, event),
        clickAddAttribute: () => clickAddAttribute(dispatch),
        attributeOnChange: (index, update) => updateAttribute(dispatch, index, update),
        removeAttribute: (index) => removeAttribute(dispatch, index)
    };
}

const ObjectEditorContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectEditor);

export default ObjectEditorContainer;
