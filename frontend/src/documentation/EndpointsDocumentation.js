import React from 'react';
import Table from 'react-bootstrap/Table';
import Required from './Required.js';

const NO_ENDPOINTS_MESSAGE = "This project doesn't have any endpoints yet.";

// ParameterRow returns the JSX for a single row of a parameter table. props should have `name`, `required`,
// `location` and `description` fields.
const ParameterRow = props => (
  <tr>
    <td>
      {props.name}
      {props.required && <Required />}
    </td>
    <td>{props.location}</td>
    <td>{props.description}</td>
  </tr>
)

// ParameterTable returns the JSX for a table of endpoint parameters. props should have a `children` field,
// containing a list of ParameterRow objects.
const ParameterTable = props => (
  <Table bordered striped className='mb-5'>
    <thead>
      <tr>
        <th>Name</th>
        <th>Location</th>
        <th>Description</th>
      </tr>
    </thead>
    <tbody>
      {
        props.children
      }
    </tbody>
  </Table>
)

// CreateDocumentation returns the JSX for the default CreateObject endpoint. props should have an `object` field
// with the object definition.
const CreateDocumentation = props => (
  <div>
    <h4>Create {props.object.name}</h4>
    <p>POST /{props.object.name.toLowerCase()}</p>
    <p>This endpoint adds a new {props.object.name} to the database.</p>
    <h5>Parameters</h5>
    <ParameterTable>
      {
        props.object.attributes.map(attribute => (
          <ParameterRow key={attribute.name} name={attribute.name} required={attribute.required} description={attribute.description} location="Body"/>
        ))
      }
    </ParameterTable>
  </div>
)

// ReadDocumentation returns the JSX for the default ReadObjectById endpoint. props should have an `object` field
// with the object definition.
const ReadDocumentation = props => {
  const idName = props.object.name.toLowerCase() + "Id"
  return (
    <div>
      <h4>Read {props.object.name} by id</h4>
      <p>GET /{props.object.name.toLowerCase()}/{"{" + idName + "}"}</p>
      <p>This endpoint returns an existing {props.object.name} from the database.</p>
      <h5>Parameters</h5>
      <ParameterTable>
        <ParameterRow name={idName} required={true} description={"The id of the " + props.object.name.toLowerCase() + " to return."} location="Path"/>
      </ParameterTable>
    </div>
  )
}

// UpdateDocumentation returns the JSX for the default UpdateObjectById endpoint. props should have an `object` field
// with the object definition.
const UpdateDocumentation = props => {
  const idName = props.object.name.toLowerCase() + "Id"
  return (
    <div>
      <h4>Update {props.object.name} by id</h4>
      <p>PUT /{props.object.name.toLowerCase()}/{"{" + idName + "}"}</p>
      <p>This endpoint updates an existing {props.object.name} in the database with the given body parameters.</p>
      <h5>Parameters</h5>
      <ParameterTable>
          <ParameterRow name={idName} required={true} description={"The id of the " + props.object.name.toLowerCase() + " to update."} location="Path"/>
          {
            props.object.attributes.map(attribute => (
              <ParameterRow key={attribute.name} name={attribute.name} required={attribute.required} description={attribute.description} location="Body"/>
            ))
          }
      </ParameterTable>
    </div>
  )
}

// DeleteDocumentation returns the JSX for the default DeleteObjectById endpoint. props should have an `object` field
// with the object definition.
const DeleteDocumentation = props => {
  const idName = props.object.name.toLowerCase() + "Id"
  return (
    <div>
      <h4>Delete {props.object.name} by id</h4>
      <p>DELETE /{props.object.name.toLowerCase()}/{"{" + idName + "}"}</p>
      <p>This endpoint removes the {props.object.name} with the given id from the database.</p>
      <h5>Parameters</h5>
      <ParameterTable>
        <ParameterRow name={idName} required={true} description={"The id of the " + props.object.name.toLowerCase() + " to delete."} location="Path"/>
      </ParameterTable>
    </div>
  )
}

// EndpointsDocumentation returns the JSX for the documentation of all endpoints. props should have an `objects` field,
// containing a list of object definition.
const EndpointsDocumentation = props => {
  if (props.objects === undefined) {
    return null;
  }

  return (
    <div>
      <h3 className="mb-0">Endpoints</h3>
      <hr className="mt-1"/>
      {
        Object.entries(props.objects).length === 0 
          ? <p>{NO_ENDPOINTS_MESSAGE}</p> 
          : Object.entries(props.objects).map(([id, object]) => (
            <div key={id}>
              <CreateDocumentation object={object} />
              <ReadDocumentation object={object} />
              <UpdateDocumentation object={object} />
              <DeleteDocumentation object={object} />
            </div>
          ))
      }
    </div>
  )
}

export default EndpointsDocumentation;
