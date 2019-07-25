import { connect } from 'react-redux';
import ObjectEditor  from '../presentation/ObjectEditor';
import {
    objectEditorUpdateName,
    objectEditorUpdateDescription,
    objectEditorAddAttribute,
    objectEditorUpdateAttribute,
    objectEditorRemoveAttribute,
    objectEditorInvalidateAttribute,
    objectEditorInvalidateDetails
}  from '../../actions/actions';

const testProps = {
    projectName: 'Test Project'
}

const INVALID_NAME_REGEX = /[^a-z ]/gi;
const REQUIRED_FIELD = "This field is required.";
const INVALID_NAME = "Only a-z, A-Z and spaces are allowed for this field.";
const INTEGER_REQUIRED = "Please specify an integer.";

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

function validateDetails(dispatch, editor) {
    let nameFeedback = null;
    if (editor.name.length === 0) {
        nameFeedback = REQUIRED_FIELD;
    } else if (editor.name.match(INVALID_NAME_REGEX) !== null) {
        nameFeedback = INVALID_NAME;
    }
    
    dispatch(objectEditorInvalidateDetails({
        'nameFeedback': nameFeedback
    }));
}

function validateObject(dispatch, editor) {
    validateDetails(dispatch, editor);
    editor.attributes.forEach(function(attribute, index) {
        let nameFeedback = null;
        if (attribute.name === null || attribute.name === '') {
            nameFeedback = REQUIRED_FIELD;
        } else if (attribute.name.match(INVALID_NAME_REGEX) !== null) {
            nameFeedback = INVALID_NAME;
        }

        let typeFeedback = null;
        if (attribute.type === null || attribute.type === 'Choose...') {
            typeFeedback = REQUIRED_FIELD;
        }

        let defaultFeedback = null;
        if (attribute.type === 'Integer' && attribute.default.length > 0) {
            const defaultValue = Number(attribute.default.replace(/,/g, ''));
            if (Number.isNaN(defaultValue) || !Number.isInteger(defaultValue)) {
                defaultFeedback = INTEGER_REQUIRED;
            }
        }

        dispatch(objectEditorInvalidateAttribute(index, {
            'nameFeedback': nameFeedback,
            'typeFeedback': typeFeedback,
            'defaultFeedback': defaultFeedback
        }));
    });
}

const mapDispatchToProps = dispatch => {
    return {
        nameOnChange: event => nameOnChange(dispatch, event),
        descriptionOnChange: event => descriptionOnChange(dispatch, event),
        clickAddAttribute: () => clickAddAttribute(dispatch),
        attributeOnChange: (index, update) => updateAttribute(dispatch, index, update),
        removeAttribute: (index) => removeAttribute(dispatch, index),
        onSubmit: (editor) => validateObject(dispatch, editor)
    };
}

const ObjectEditorContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectEditor);

export default ObjectEditorContainer;
