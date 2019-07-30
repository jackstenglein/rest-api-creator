import React from 'react';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Breadcrumb from 'react-bootstrap/Breadcrumb';
import IconButton from '@material-ui/core/IconButton';
import Delete from '@material-ui/icons/Delete';
import {
    OBJECT_EDITOR_EDITING,
    OBJECT_EDITOR_REQUEST_PENDING,
    OBJECT_EDITOR_REQUEST_SUCCEEDED,
    OBJECT_EDITOR_REQUEST_FAILED,
    OBJECT_EDITOR_CANCEL
} from '../../redux/actions/objects/objectEditorActions.js';
import { Redirect } from 'react-router-dom';
import Modal from 'react-bootstrap/Modal';
import Spinner from 'react-bootstrap/Spinner';


const OBJECT_DESC_PLACEHOLDER = "Add a description of what this object represents"
    + " and it will be included in your auto-generated documentation.";
const ATTR_DESC_PLACEHOLDER = "Add a description of what this attribute represents"
    + " and it will be included in your auto-generated documentation.";
const ID_DESCRIPTION = "The id attribute uniquely identifies database records of this"
    + " object type. There can only be one record with a given id. Each time an"
    + " object of this type is created, it is automatically given a new id.";

function getAttributeRowJsx(index, attribute, onChange) {
    return (
        <Form.Row key={"attributeRow" + index}>
            <Form.Group as={Col}>
                <Form.Label>Attribute name</Form.Label>
                <Form.Control
                    value={attribute.name}
                    onChange={(event) => onChange(index, {"name": event.target.value})}
                    placeholder="Enter name"
                    isInvalid={attribute.nameFeedback !== null}
                />
                <Form.Control.Feedback type="invalid">{attribute.nameFeedback}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group as={Col}>
                <Form.Label>Data type</Form.Label>
                <Form.Control
                    as="select"
                    value={attribute.type}
                    onChange={(event) => onChange(index, {"type": event.target.value})}
                    isInvalid={attribute.typeFeedback !== null}
                >
                    <option>Choose...</option>
                    <option>Text</option>
                    <option>Integer</option>
                </Form.Control>
                <Form.Control.Feedback type="invalid">{attribute.typeFeedback}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group as={Col}>
                <Form.Label>Default value</Form.Label>
                <Form.Control
                    value={attribute.default}
                    onChange={(event) => onChange(index, {"default": event.target.value})}
                    placeholder="Optional value"
                    isInvalid={attribute.defaultFeedback !== null}
                />
                <Form.Control.Feedback type="invalid">{attribute.defaultFeedback}</Form.Control.Feedback>
            </Form.Group>
        </Form.Row>
    );
}

function getAttributeDescriptionJsx(index, attribute, onChange) {
    return (
        <Form.Group key={"attributeDescription" + index}>
            <Form.Label>Description</Form.Label>
            <Form.Control
                as="textarea"
                rows="2"
                value={attribute.description}
                onChange={(event) => onChange(index, {"description": event.target.value})}
                placeholder={ATTR_DESC_PLACEHOLDER}
            >
            </Form.Control>
        </Form.Group>
    )
}

function getAttribute(index, attribute, onChange, onClickRemove) {
    return (
        <Row key={"attribute" + index} className="align-items-center">
            <Col xs="11">
                { getAttributeRowJsx(index, attribute, onChange) }
                { getAttributeDescriptionJsx(index, attribute, onChange) }
            </Col>
            <Col xs="1">
                <IconButton onClick={() => onClickRemove(index)}><Delete/></IconButton>
            </Col>
        </Row>
    );
}

function displayAttributes(attributes, callbacks) {
    let attributeJsx = [];
    for (let i = 0; i < attributes.length; ++i) {
        attributeJsx.push(<hr key={"separator" + i}/>);
        attributeJsx.push(getAttribute(i, attributes[i], callbacks.attributeOnChange, callbacks.removeAttribute));
    }
    return attributeJsx;
}

function getBreadcrumbJsx(props) {
    if (props.control.selectedObject === -1) {
        return (
            <Breadcrumb bsPrefix="breadcrumb crud">
                <Breadcrumb.Item href="#">Projects</Breadcrumb.Item>
                <Breadcrumb.Item href="#">{props.projectName}</Breadcrumb.Item>
                <Breadcrumb.Item href="#">Objects</Breadcrumb.Item>
                <Breadcrumb.Item active>Create</Breadcrumb.Item>
            </Breadcrumb>
        );
    } else {
        return (
            <Breadcrumb bsPrefix="breadcrumb crud">
                <Breadcrumb.Item href="#">Projects</Breadcrumb.Item>
                <Breadcrumb.Item href="#">{props.projectName}</Breadcrumb.Item>
                <Breadcrumb.Item href="#">Objects</Breadcrumb.Item>
                <Breadcrumb.Item active>{props.details.name}</Breadcrumb.Item>
                <Breadcrumb.Item active>Edit</Breadcrumb.Item>
            </Breadcrumb>
        );
    }
}

function displayIdAttribute() {
    return (
        <>
            <Form.Row key={"attributeRowId"}>
                <Form.Group as={Col}>
                    <Form.Label>Attribute name</Form.Label>
                    <Form.Control value="id" disabled />
                </Form.Group>
                <Form.Group as={Col}>
                    <Form.Label>Data type</Form.Label>
                    <Form.Control value="Integer" disabled />
                </Form.Group>
                <Form.Group as={Col}>
                    <Form.Label>Default value</Form.Label>
                    <Form.Control value="Auto-Increment" disabled />
                </Form.Group>
            </Form.Row>
            <Form.Group key={"attributeDescriptionId"}>
                <Form.Label>Description</Form.Label>
                <Form.Control as="textarea" rows="2" value={ID_DESCRIPTION} disabled />
            </Form.Group>
        </>
    );
}

function getTopBar(props) {
    let buttonTitle;
    let buttonsDisabled = false;
    if (props.control.selectedObject === -1) {
        buttonTitle = "Create";
    } else {
        buttonTitle = "Save";
    }

    let submitButton;
    let cancelButton;
    if (props.control.status === OBJECT_EDITOR_REQUEST_PENDING) {
        submitButton = (
            <Button variant="primary" className="mr-2" disabled>
                <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true"/>
            </Button>
        );
        cancelButton = (<Button variant="danger" disabled>Cancel</Button>);
    } else {
        submitButton = (
            <Button variant="primary" className="mr-2" onClick={() => props.callbacks.onSubmit(props)}>{buttonTitle}</Button>
        );
        cancelButton = (<Button variant="danger" onClick={() => props.callbacks.cancel()}>Cancel</Button>);
    }

    return (
        <Row className="align-items-center justify-content-between object-editor-toolbar">
            <Col xs="auto">
                { getBreadcrumbJsx(props) }
            </Col>
            <Col xs="auto">
                <ButtonToolbar>
                    { submitButton }
                    { cancelButton }
                </ButtonToolbar>
            </Col>
        </Row>
    );
}

function getErrorModal(control, callbacks) {
    if (control.status !== OBJECT_EDITOR_REQUEST_FAILED) {
        return null;
    }

    return (
        <Modal centered show onHide={callbacks.closeErrorModal}>
            <Modal.Header closeButton>
                <Modal.Title>Unable to create object</Modal.Title>
            </Modal.Header>

            <Modal.Body>
                <p>{control.errorMessage}</p>
            </Modal.Body>
        </Modal>
    );
}

function ObjectEditor(props) {
    const { control, details, attributes, callbacks } = props;

    if (control.status === OBJECT_EDITOR_REQUEST_SUCCEEDED ||
        control.status === OBJECT_EDITOR_CANCEL) {
        return (<Redirect to="/objects" />);
    }

    return (
        <Container className="object-editor-viewport">
            { getErrorModal(control, callbacks) }
            { getTopBar(props) }
            <h4>Details</h4>
            <Form.Group controlId="createObjectName">
                <Form.Label>Object name</Form.Label>
                <Form.Control
                    onChange={callbacks.nameOnChange}
                    value={details.name}
                    placeholder="Enter name"
                    isInvalid={details.nameFeedback !== null}
                />
                <Form.Control.Feedback type="invalid">{details.nameFeedback}</Form.Control.Feedback>
            </Form.Group>
            <Form.Group controlId="createObjectDescription">
                <Form.Label>Description</Form.Label>
                <Form.Control
                    as="textarea"
                    rows="3"
                    placeholder={OBJECT_DESC_PLACEHOLDER}
                    onChange={callbacks.descriptionOnChange}
                    value={details.description}
                />
            </Form.Group>
            <h4>Attributes</h4>
            { displayIdAttribute() }
            { displayAttributes(attributes, callbacks) }
            <Button variant="primary" onClick={callbacks.clickAddAttribute}>Add attribute</Button>
        </Container>
    );
}

export default ObjectEditor;
