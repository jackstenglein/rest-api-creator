import React from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import ObjectHelp from './objects/ObjectHelp';
import AttributeHelp from './objects/AttributeHelp';


export const OBJECT_HELP = "help/objects"
export const ATTRIBUTE_HELP = "help/attributes"


const Help = props => {

  var help;
  switch (props.id) {
    case OBJECT_HELP:
      help = ObjectHelp;
      break;
    case ATTRIBUTE_HELP:
      help = AttributeHelp;
      break;
  }

  if (help === undefined) {
    return null;
  }

  return (
    <Col className="border-left" xs={4}>
      <Row className="mb-2">
        <Col>
          <h4>{help.title}</h4>
        </Col>
        <Col>
          <button type="button" className="close" aria-label="Close" onClick={props.close}>
            <span aria-hidden="true">&times;</span>
          </button>
        </Col>
      </Row>
      { help.paragraphs.map((paragraph, idx) => (
        <p key={idx}>{paragraph}</p>
      ))}
    </Col>
  )
}

export default Help;
