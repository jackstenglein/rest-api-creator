import React from 'react';
import { connect } from 'react-redux';
import { putObjectSuccess } from '../../redux/modules/projects'; 
import ObjectEditor from './ObjectEditor.js';
import produce from 'immer';
import validateObject from './Validator.js';
import deepEqual from 'deep-equal';
import { Prompt } from 'react-router-dom';

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
   

    console.log("Props: ", props)

    this.state = {
      values: this.getOriginalObject()
    };

    this.addAttribute = this.addAttribute.bind(this);
    this.changeName = this.changeName.bind(this);
    this.changeDescription = this.changeDescription.bind(this);
    this.changeAttribute = this.changeAttribute.bind(this);
    this.onSave = this.onSave.bind(this);
    this.onCancel = this.onCancel.bind(this);
    this.removeAttribute = this.removeAttribute.bind(this);
  }

  // addAttribute adds an empty attribute object to the state's values.attributes array. It also adds an empty 
  // error object to the state's errors.attributes array.
  addAttribute() {
    const nextState = produce(this.state, draftState => {
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
    const object = this.props.allObjects[objectId];
    if (object !== undefined) {
      return JSON.parse(JSON.stringify(object));
    } 
    return JSON.parse(JSON.stringify(EMPTY_OBJECT));
  }

  // onSave handles clicks to the save button. If this action is triggered, then the object should be valid.
  onSave() {
    this.props.onSave(this.state.values);
  }

  // onCancel handles clicks to the cancel button.
  onCancel() {
    console.log("Cancel");
  }

  // removeAttribute removes the attribute at the given index from the state's values.attributes array.
  // removeAttribute also removes the error object at the given index from the state's errors.attributes array.
  removeAttribute(i) {
    const nextState = produce(this.state, draftState => {
      draftState.values.attributes.splice(i, 1);
      draftState.errors.attributes.splice(i, 1);
    })
    this.setState(nextState);
  }

  // render returns the JSX for the ObjectEditor.
  render() {
    const onChangeHandlers = {
      name: this.changeName,
      description: this.changeDescription,
      attribute: this.changeAttribute
    }

    const [isValid, errors] = validateObject(this.state.values, this.props.allObjects);
    const cancelPrompt = !deepEqual(this.state.values, this.getOriginalObject())
    console.log("Cancel prompt: ", cancelPrompt)

    return (
      <>
        <Prompt when={cancelPrompt} message="Are you sure you want to discard your changes?"/> 
        <ObjectEditor 
          values={this.state.values}
          isValid={isValid}
          errors={errors} 
          onChangeHandlers={onChangeHandlers}
          onSave={this.onSave}
          onCancel={this.onCancel}
          removeAttribute={this.removeAttribute}
          addAttribute={this.addAttribute}
        />
      </>
    )
  }
}

const mapStateToProps = state => {
  return {
    allObjects: state.projects.defaultProject.objects
  };
}

const mapDispatchToProps = dispatch => {
  return {
    onSave: object => {
      dispatch(putObjectSuccess("defaultProject", object))
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(ObjectEditorContainer);
