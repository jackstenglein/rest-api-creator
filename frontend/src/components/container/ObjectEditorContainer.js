import { connect } from 'react-redux';
import ObjectEditor  from '../presentation/ObjectEditor';
import {
    objectEditorUpdateName,
    objectEditorUpdateDescription,
    objectEditorAddAttribute
}  from '../../actions/actions';

const mapStateToProps = (state, ownProps) => {
    return Object.assign({}, state.objects.editor);
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

const mapDispatchToProps = dispatch => {
    return {
        nameOnChange: event => nameOnChange(dispatch, event),
        descriptionOnChange: event => descriptionOnChange(dispatch, event),
        clickAddAttribute: () => clickAddAttribute(dispatch)
    };
}

const ObjectEditorContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectEditor);

export default ObjectEditorContainer;
