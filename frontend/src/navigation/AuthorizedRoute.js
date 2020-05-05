import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';
import * as network from '../redux/modules/network.js';
import * as userInfo from '../redux/modules/userInfo.js';
import { getUser } from '../api/api.js';

class AuthorizedRoute extends React.Component {
  async componentDidMount() {
    const authenticated = this.props.userInfo.authenticated;
    const net = this.props.userInfo.network;
    if (!authenticated && (net.status === network.STATUS_NONE || net.status === network.STATUS_FAILURE)) {
      console.log("Query for the user");
      this.props.getUserRequest()
      const response = await getUser();
      this.props.getUserResponse(response);
    }
  }

  render() {
    const { component: Component, userInfo, shouldAuthenticate, ...rest} = this.props;
    const status = userInfo.network.status;

    return (
      <Route {...rest} render={({location}) => {
        // The GetUser request completed
        if (status === network.STATUS_SUCCESS || status === network.STATUS_FAILURE) {
          // We have the desired authentication level, so return the component.
          if (userInfo.authenticated === shouldAuthenticate) {
            return <Component/>
          }
          // We are either not logged in and should be, or we are logged in and shouldn't be. Redirect as appropriate
          return shouldAuthenticate ? <Redirect to={{pathname: "/login", state: {from: location}}}/> 
            : <Redirect to="/app/details"/>;
        }

        // The GetUser request is still loading
        if (status === network.STATUS_PENDING || status === network.STATUS_NONE) {
          return null;
        }

        // We should never get here
        return <p>If you are seeing this, there is a bug in my code. Please let me know. Thanks.</p>  
      }}/>
    )
  }
}

const mapStateToProps = state => {
  return {
    userInfo: state.userInfo
  };
}

const mapDispatchToProps = dispatch => {
  return {
    getUserRequest: () => {
      dispatch(userInfo.getUserRequest())
    },
    getUserResponse: response => {
      dispatch(userInfo.getUserResponse(response))
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(AuthorizedRoute);
