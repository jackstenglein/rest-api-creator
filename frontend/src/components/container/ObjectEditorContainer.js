import { connect } from 'react-redux';
import ObjectEditor  from '../presentation/ObjectEditor';
import {
    objectEditorUpdateName,
    objectEditorUpdateDescription,
    objectEditorAddAttribute,
    objectEditorUpdateAttribute,
    objectEditorRemoveAttribute,
    objectEditorInvalidateAttribute,
    objectEditorInvalidateDetails,
    objectEditorCloseErrorModal,
    objectEditorCancel,
    createObject
} from '../../redux/actions/objects/objectEditorActions.js';

const testProps = {
    projectName: 'Test Project'
}

const INVALID_NAME_REGEX = /[^a-z ]/gi;
const REQUIRED_FIELD = "This field is required.";
const INVALID_NAME = "Only a-z, A-Z and spaces are allowed for this field.";
const INTEGER_REQUIRED = "Please specify an integer.";
const DUPLICATE_ATTRIBUTE  = "An attribute with this name already exists for this object.";

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

    return nameFeedback === null;
}

function validateObject(dispatch, editor) {
    let detailsValidated = validateDetails(dispatch, editor.details);
    let attributesValidated = true;
    let attributeNames = new Set([]);
    editor.attributes.forEach(function(attribute, index) {
        let nameFeedback = null;
        if (attribute.name === null || attribute.name === '') {
            nameFeedback = REQUIRED_FIELD;
        } else if (attribute.name.match(INVALID_NAME_REGEX) !== null) {
            nameFeedback = INVALID_NAME;
        } else if (attributeNames.has(attribute.name)) {
            nameFeedback = DUPLICATE_ATTRIBUTE;
        } else {
            attributeNames.add(attribute.name);
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

        attributesValidated = (attributesValidated && nameFeedback === null
            && typeFeedback === null && defaultFeedback === null);

        dispatch(objectEditorInvalidateAttribute(index, {
            'nameFeedback': nameFeedback,
            'typeFeedback': typeFeedback,
            'defaultFeedback': defaultFeedback
        }));
    });

    return detailsValidated && attributesValidated;
}

function submitObject(dispatch, editor) {
    let objectValidated = validateObject(dispatch, editor);
    if (objectValidated) {
        console.log("OBJECT VALIDATED");
        dispatch(createObject({name: editor.details.name, attributes: editor.attributes}));
    } else {
        console.log("OBJECT INVALID");
    }
}

const mapDispatchToProps = dispatch => {
    return {
        callbacks: {
            nameOnChange: event => nameOnChange(dispatch, event),
            descriptionOnChange: event => descriptionOnChange(dispatch, event),
            clickAddAttribute: () => clickAddAttribute(dispatch),
            attributeOnChange: (index, update) => updateAttribute(dispatch, index, update),
            removeAttribute: (index) => removeAttribute(dispatch, index),
            onSubmit: (editor) => submitObject(dispatch, editor),
            closeErrorModal: () => dispatch(objectEditorCloseErrorModal()),
            cancel: () => dispatch(objectEditorCancel())
        }
    };
}

const ObjectEditorContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectEditor);

export default ObjectEditorContainer;
