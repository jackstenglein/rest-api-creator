import React from 'react';
import { connect } from 'react-redux';
import { putObjectSuccess } from '../../redux/modules/projects'; 
import ObjectEditor from './ObjectEditor.js';
import produce from 'immer';
import validateObject from './Validator.js';
import deepEqual from 'deep-equal';
import { Prompt, Redirect } from 'react-router-dom';
import { putObject } from '../../api/api.js';
import { STATUS_SUCCESS } from '../../redux/modules/network.js';

const EMPTY_OBJECT = {
  name: "",
  description: "",
  attributes: [
    {
      name: "",
      type: "",
      defaultValue: "",
      description: ""
    }
  ]
}

class ObjectEditorContainer extends React.Component {

  constructor(props) {
    super(props)
   
    this.state = {values: this.getOriginalObject()};
    this.addAttribute = this.addAttribute.bind(this);
    this.changeName = this.changeName.bind(this);
    this.changeDescription = this.changeDescription.bind(this);
    this.changeAttribute = this.changeAttribute.bind(this);
    this.onSave = this.onSave.bind(this);
    this.removeAttribute = this.removeAttribute.bind(this);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.project.network.status !== STATUS_SUCCESS && this.props.project.network.status === STATUS_SUCCESS) {
      this.setState({values: this.getOriginalObject()});
    }
  }

  // addAttribute adds an empty attribute object to the state's values.attributes array.
  addAttribute() {
    const nextState = produce(this.state, draftState => {
      if (draftState.values.attributes === undefined) {
        draftState.values.attributes = [];
      }
      draftState.values.attributes.push({name: "", type: "", defaultValue: "", description: ""});
    })
    this.setState(nextState);
  }

  // changeName handles input to the object's name field.
  changeName(event) {
    const nextState = produce(this.state, draftState => {
      draftState.values.name = event.target.value;
    })
    this.setState(nextState);
  }

  // changeDescription handles input to the object's description field.
  changeDescription(event) {
    const nextState = produce(this.state, draftState => {
      draftState.values.description = event.target.value;
    })
    this.setState(nextState);
  }

  // changeAttribute handles input to a specific attribute's fields. changeAttribute receives the index of the 
  // attribute, the field to change and the input event.
  changeAttribute(i, field, event) {
    const nextState = produce(this.state, draftState => {
      draftState.values.attributes[i][field] = event.target.value;
    })
    this.setState(nextState);
  }

  // getOriginalObject returns the starting object definition for when the object editor first opens.
  getOriginalObject() {
    const objectId = this.props.match.params.objectId;
    if (objectId === undefined) {
      // We are creating an object
      return JSON.parse(JSON.stringify(EMPTY_OBJECT));
    } 

    if (this.props.project.network.status !== STATUS_SUCCESS) {
      // We haven't finished fetching the project
      return undefined;
    }

    const object = this.props.project.objects[objectId];
    if (object !== undefined) {
      // We are editing an existing object
      return JSON.parse(JSON.stringify(object));
    } 

    // We are trying to edit an object but it does not exist
    return undefined;
  }

  // onSave handles clicks to the save button. If this action is triggered, then the object should be valid.
  async onSave() {
    // Make the request
    console.log("Making putObject request")
    const response = await putObject("defaultProject", this.state.values);
    console.log("Got response: ", response);

    // Make changes to state based on response
    var apiError = "";
    var saved = false;
    if (response === undefined) {
      apiError = "Failed to make network request.";
    } else if (response.error !== undefined) {
      apiError = response.error;
    } else {
      saved = true;
      this.props.onSave(this.state.values);
    }

    const nextState = produce(this.state, draftState => {
      draftState.apiError = apiError;
      draftState.saved = saved;
    });
    this.setState(nextState);
  }

  // removeAttribute removes the attribute at the given index from the state's values.attributes array.
  removeAttribute(i) {
    const nextState = produce(this.state, draftState => {
      draftState.values.attributes.splice(i, 1);
    })
    this.setState(nextState);
  }

  // render returns the JSX for the ObjectEditor.
  render() {
    if (this.state.saved) {
      return <Redirect to="/app/objects"/>
    }

    const onChangeHandlers = {
      name: this.changeName,
      description: this.changeDescription,
      attribute: this.changeAttribute
    }

    const [isValid, errors] = validateObject(this.state.values, this.props.project.objects);
    const cancelPrompt = !deepEqual(this.state.values, this.getOriginalObject())
    const alertError = this.state.apiError || this.props.project.network.error;

    return (
      <>
        <Prompt when={cancelPrompt} message="Are you sure you want to discard your changes?"/> 
        <ObjectEditor 
          values={this.state.values}
          alertError={alertError}
          isValid={isValid}
          errors={errors} 
          onChangeHandlers={onChangeHandlers}
          onSave={this.onSave}
          removeAttribute={this.removeAttribute}
          addAttribute={this.addAttribute}
        />
      </>
    )
  }
}

const mapDispatchToProps = dispatch => {
  return {
    onSave: object => {
      dispatch(putObjectSuccess("defaultProject", object))
    }
  }
}

export default connect(null, mapDispatchToProps)(ObjectEditorContainer);
