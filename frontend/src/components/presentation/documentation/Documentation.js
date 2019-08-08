import React from 'react';
import Container from 'react-bootstrap/Container';
import ObjectDocumentation from './ObjectDocumentation';

const testObject = {
    id: 31,
    details: {
        'name': 'Test Object Name',
        'description': 'This is a test object description. It is hopefully long enough to span across two lines.'
    },
    attributes: [
        {
            'name': 'name',
            'type': 'Text',
            'description': 'The name of the user.'
        }
    ]
};

const testObject2 = {
    id: 31,
    details: {
        'name': 'Test Object 2 Name',
        'description': 'This is a test object description. It is hopefully long enough to span across two lines.'
    },
    attributes: [
        {
            'name': 'name2',
            'type': 'Text',
            'description': 'The name of the user.'
        }
    ]
};

const testObject3 = {
    id: 31,
    details: {
        'name': 'Test Object 3 Name',
        'description': 'This is a test object description. It is hopefully long enough to span across two lines.'
    },
    attributes: [
        {
            'name': 'name2',
            'type': 'Text',
            'description': 'The name of the user.'
        }
    ]
};

const testObjects = [testObject, testObject2, testObject3];

function Documentation(props) {
    return (
        <Container>
            <h2 className='mb-0 mt-3'>Objects</h2>
            <hr className='mt-0'/>
            {
                testObjects.map(function(object) {
                    return (<ObjectDocumentation key={object.id} object={object}/>);
                })
            }
        </Container>
    );
}

export default Documentation;
