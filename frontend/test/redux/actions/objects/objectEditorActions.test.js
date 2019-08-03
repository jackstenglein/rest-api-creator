import * as Actions from '../../../../src/redux/actions/objects/objectEditorActions.js';

test('objectEditorUpdateName action creator', () => {
    expect(Actions.objectEditorUpdateName('Test name')).toEqual({
        type: 'ObjectEditorUpdateName',
        details: {
            name: 'Test name'
        }
    });
});

test('objectEditorUpdateDescription action creator', () => {
    expect(Actions.objectEditorUpdateDescription('Test description')).toEqual({
        type: 'ObjectEditorUpdateDescription',
        details: {
            description: 'Test description'
        }
    });
});

test('objectEditorAddAttribute action creator', () => {
    expect(Actions.objectEditorAddAttribute()).toEqual({type: 'ObjectEditorAddAttribute'});
});

test('objectEditorUpdateAttribute action creator', () => {
    const index = 5;
    const attribute = {
        name: 'Test attribute',
        description: 'Test attribute description',
        type: 'Text'
    };
    expect(Actions.objectEditorUpdateAttribute(index, attribute)).toEqual({
        type: 'ObjectEditorUpdateAttribute',
        index: index,
        update: attribute
    });
});

test('objectEditorRemoveAttribute action creator', () => {
    const index = 7;
    expect(Actions.objectEditorRemoveAttribute(index)).toEqual({type: 'ObjectEditorRemoveAttribute', index: index});
});

test('objectEditorInvalidateAttribute action creator', () => {
    const index = 8;
    const feedback = {
        nameFeedback: 'Incorrect name.',
        typeFeedback: 'Incorrect type.',
        defaultFeedback: 'Incorrect default'
    };
    expect(Actions.objectEditorInvalidateAttribute(index, feedback)).toEqual({
        type: 'ObjectEditorInvalidateAttribute',
        index: index,
        feedback: feedback
    });
});

test('objectEditorInvalidateDetails action creator', () => {
    const feedback = {
        nameFeedback: 'Name is not good'
    };
    expect(Actions.objectEditorInvalidateDetails(feedback)).toEqual({
        type: 'ObjectEditorInvalidateDetails',
        details: {
            nameFeedback: feedback.nameFeedback
        }
    });
});

test('objectEditorCloseErrorModal action creator', () => {
    expect(Actions.objectEditorCloseErrorModal()).toEqual({type: 'ObjectEditorCloseErrorModal'});
});

test('objectEditorCancel action creator', () => {
    expect(Actions.objectEditorCancel()).toEqual({type: 'ObjectEditorCancel'});
});
