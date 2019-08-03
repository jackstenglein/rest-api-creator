import * as Actions from '../../../../src/redux/actions/objects/objectViewActions.js';

test('clickEdit action creator', () => {
    const object = {
        name: 'Test object',
        description: 'Test description',
        attributes: [
            {
                name: 'Test attribute',
                description: 'Test attribute description',
                type: 'Integer'
            }
        ]
    };
    expect(Actions.clickEdit(object)).toEqual({type: 'ObjectViewClickEdit', object: object})
});
