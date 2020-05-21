import React from 'react';
import ObjectAttributeEditor from './ObjectAttributeEditor';
import ObjectInfoEditor from './ObjectInfoEditor';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';
import ErrorAlert from '../../portal/ErrorAlert.js';
import Row from 'react-bootstrap/Row';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import { Link } from 'react-router-dom';
import Help, { OBJECT_HELP } from '../../help/Help.js';

const EditorToolbar = props => (
  <Row className="align-items-center justify-content-between">
    <Col xs="auto">
      <h2>
        Edit Object <i className="material-icons clickable" onClick={() => props.toggleHelp(OBJECT_HELP)}>help_outline</i>
      </h2>
    </Col>
    <Col xs="auto">
      <ButtonToolbar>
        <Button className="mr-2" onClick={props.onSave} disabled={!props.isValid}>Save</Button>
        <Link to="/app/objects"><Button variant="danger">Cancel</Button></Link>
      </ButtonToolbar>
    </Col>
  </Row>
)

const ObjectEditor = props => {
  if (props.values === undefined) {
    return null;
  }

  return (
  <div>
    <Row>
    <Col>
    <EditorToolbar isValid={props.isValid} onSave={props.onSave} onCancel={props.onCancel} toggleHelp={props.toggleHelp}/>
    <br />
    <ErrorAlert error={props.alertError}/>
    
    <ObjectInfoEditor {...props}/>
    <ObjectAttributeEditor 
      values={props.values.attributes} 
      errors={props.errors.attributes} 
      onChange={props.onChangeHandlers.attribute}
      remove={props.removeAttribute}
      add={props.addAttribute}
      toggleHelp={props.toggleHelp}
    />
    </Col>
    {props.helpId && (
      
        <Help id={props.helpId} close={props.toggleHelp}/>
    )}
   </Row>
  </div>
)
}

export default ObjectEditor;
