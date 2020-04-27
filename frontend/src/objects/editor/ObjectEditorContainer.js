import React from 'react';
import ObjectEditor from './ObjectEditor.js';
import produce from 'immer';
import validateObject from './Validator.js';

const EMPTY_OBJECT = {
  name: "",
  description: "",
  attributes: [
    {
      name: "",
      type: "",
      defaultValue: "",
      description: ""
    },
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
    var originalObject;
    if (props.object !== undefined) {
      originalObject = props.object;
    } else {
      originalObject = EMPTY_OBJECT;
    }

    this.state = {
      values: JSON.parse(JSON.stringify(originalObject)),
      errors: EMPTY_OBJECT
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
      draftState.errors.attributes.push({name: "", type: "", defaultValue: ""});
    })
    this.setState(nextState);
  }

  // changeName handles input to the object's name field.
  changeName(event) {
    const nextState = produce(this.state, draftState => {
      draftState.values.name = event.target.value;
      draftState.errors = validateObject(draftState.values);
    })
    this.setState(nextState);
  }

  // changeDescription handles input to the object's description field.
  changeDescription(event) {
    const nextState = produce(this.state, draftState => {
      draftState.values.description = event.target.value;
      draftState.errors = validateObject(draftState.values);
    })
    this.setState(nextState);
  }

  // changeAttribute handles input to a specific attribute's fields. changeAttribute receives the index of the 
  // attribute, the field to change and the input event.
  changeAttribute(i, field, event) {
    const nextState = produce(this.state, draftState => {
      draftState.values.attributes[i][field] = event.target.value;
      draftState.errors = validateObject(draftState.values);
    })
    this.setState(nextState);
  }

  // onSave handles clicks to the save button.
  onSave() {
    console.log("Saved object: ", this.state.values);
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

    return (
      <ObjectEditor 
        values={this.state.values}
        errors={this.state.errors} 
        onChangeHandlers={onChangeHandlers}
        onSave={this.onSave}
        onCancel={this.onCancel}
        removeAttribute={this.removeAttribute}
        addAttribute={this.addAttribute}
      />
    )
  }
}

export default ObjectEditorContainer;
