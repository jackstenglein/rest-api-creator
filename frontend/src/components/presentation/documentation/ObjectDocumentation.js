import React from 'react';
import Table from 'react-bootstrap/Table';

const ID_DESCRIPTION = "The id attribute uniquely identifies database records of this"
    + " object type. There can only be one record with a given id. Each time an"
    + " object of this type is created, it is automatically given a new id.";

function getIdAttributeJsx() {
    return (
        <tr>
            <td>id</td>
            <td>Integer</td>
            <td>{ID_DESCRIPTION}</td>
        </tr>
    );
}

function ObjectDocumentation(props) {
    const object = props.object;
    return (
        <>
            <h5>{object.details.name}</h5>
            <p>{object.details.description}</p>
            <h6>Attributes</h6>
            <Table bordered striped className='mb-5'>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th className="text-nowrap">Data Type</th>
                        <th>Description</th>
                    </tr>
                </thead>
                <tbody>
                    { getIdAttributeJsx() }
                    {
                        object.attributes.map(function(attribute) {
                            return (
                                <tr key={attribute.name}>
                                    <td>{attribute.name}</td>
                                    <td>{attribute.type}</td>
                                    <td>{attribute.description}</td>
                                </tr>
                            )
                        })
                    }
                </tbody>
            </Table>
        </>
    );
}

export default ObjectDocumentation;
