import React from 'react';
import Table from 'react-bootstrap/Table';
import Required from './Required.js';

const ID_DESCRIPTION = "The id attribute uniquely identifies database records of this"
    + " object type. There can only be one record with a given id. Each time an"
    + " object of this type is created, it is automatically given a new id. The id can"
    + " be used to read and edit specific object instances from the database.";
const NO_OBJECTS_MESSAGE = "This project doesn't have any objects yet.";

// AttributeDocumentation returns the JSX for the documentation of a single attribute. props should have an
// `attributes` field, containing a dictionary with the attribute definition. AttributeDocumentation returns
// a <tr> element, so it should be embedded into a table.
const AttributeDocumentation = props => (
  <tr>
    <td>
      {props.attribute.name}
      {props.attribute.required && <Required />}
    </td>
    <td>{props.attribute.type}</td>
    <td>{props.attribute.description}</td>
  </tr>
)

// ObjectDocumentation returns the JSX for the documentation of a single object. props should have an
// `object` field, containing a dictionary with the object definition.
const ObjectDocumentation = props => (
  <div>
    <h4>{props.object.name}</h4>
    <p>{props.object.description}</p>
    <p>{props.object.name} has the following attributes:</p>
    <Table bordered striped className='mb-5'>
      <thead>
        <tr>
          <th>Name</th>
          <th className="text-nowrap">Data Type</th>
          <th>Description</th>
        </tr>
      </thead>
      <tbody>
        <AttributeDocumentation attribute={{name: "id", type: "Integer", description: ID_DESCRIPTION}}/>
        {
          props.object.attributes && props.object.attributes.map(attribute => (
            <AttributeDocumentation key={attribute.name} attribute={attribute} />
          ))
        }
      </tbody>
    </Table>
  </div>
)

// ObjectsDocumentation returns the JSX for the documentation of all objects. props should have an
// `objects` field, containing an array with the object definitions.
const ObjectsDocumentation = props => {
  if (props.objects === undefined) {
    return null;
  }

  return (
    <div>
      <h3 className="mb-0">Objects</h3>
      <hr className="mt-1"/>
      {
        !props.objects || Object.entries(props.objects).length === 0 
          ? <p>{NO_OBJECTS_MESSAGE}</p> 
          : Object.entries(props.objects).map(([id, object]) => (
            <ObjectDocumentation key={id} object={object} />
          ))
      }
    </div>
  )
}

export default ObjectsDocumentation;
