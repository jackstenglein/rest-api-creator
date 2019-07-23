import React from 'react';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';

const DESCRIPTION_PLACEHOLDER = "Add a description of what this object represents"
    + " and it will be included in your auto-generated documentation.";

function testOnChange(event, i) {
    console.log("On change for %d: ", i, event.target.value);
}

function displayAttributes(attributes) {
    let attributeJsx = [];
    for (let i = 0; i < attributes.length; ++i) {
        attributeJsx.push(
            <Form.Row key={i}>
                <Form.Group as={Col}>
                    <Form.Label>Attribute name</Form.Label>
                    <Form.Control placeholder="Enter name" />
                </Form.Group>
                <Form.Group as={Col}>
                    <Form.Label>Data type</Form.Label>
                    <Form.Control as="select" value="Text" onChange={(event) => testOnChange(event, i)}>
                        <option>Choose...</option>
                        <option>Text</option>
                        <option>Integer</option>
                    </Form.Control>
                </Form.Group>
                <Form.Group as={Col}>
                    <Form.Label>Default value</Form.Label>
                    <Form.Control placeholder="Optional value" />
                </Form.Group>
            </Form.Row>
        );
        attributeJsx.push(<hr />);
    }

    attributeJsx.pop();
    return attributeJsx;
}

function ObjectEditor(props) {
    return (
        <Container>
            <Form.Group controlId="createObjectName">
                <Form.Label>Object name</Form.Label>
                <Form.Control onChange={props.nameOnChange} value={props.name} placeholder="Enter name" />
            </Form.Group>
            <Form.Group controlId="createObjectDescription">
                <Form.Label>Description</Form.Label>
                <Form.Control
                    as="textarea"
                    rows="3"
                    placeholder={DESCRIPTION_PLACEHOLDER}
                    onChange={props.descriptionOnChange}
                    value={props.description}
                />
            </Form.Group>
            <h4>Attributes</h4>
            { displayAttributes(props.attributes) }
            <Button variant="primary" onClick={props.clickAddAttribute}>Add attribute</Button>
        </Container>
    );
}

export default ObjectEditor;
