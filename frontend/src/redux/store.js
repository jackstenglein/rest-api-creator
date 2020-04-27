// store.js creates the Redux store.

import { combineReducers } from 'redux';
import { configureStore } from '@reduxjs/toolkit';
import logger from 'redux-logger';
import * as reducers from './modules/index';

const rootReducer = combineReducers(reducers);
const store = configureStore({
  reducer: rootReducer,
  middleware: [logger],
  devTools: process.env.NODE_ENV !== 'production'
});

export default store;
