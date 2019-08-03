import * as Actions from '../../../../src/redux/actions/objects/objectListActions.js';

test('clickCreate action creator', () => {
    expect(Actions.clickCreate()).toEqual({type: 'ObjectListViewClickCreate'});
});

test('clickObject action creator', () => {
    expect(Actions.clickObject(5)).toEqual({type: 'ObjectListViewClickObject', id: 5});
});

test('resetListView action creator', () => {
    expect(Actions.resetListView()).toEqual({type: 'ObjectListViewReset'});
});
