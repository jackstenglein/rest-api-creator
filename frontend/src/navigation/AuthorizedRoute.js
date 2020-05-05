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
    const { component: Component, userInfo, ...rest} = this.props;
    return (
      <Route {...rest} render={({location}) => {
        if (userInfo.authenticated) {
          return <Component/>
        }

        switch(userInfo.network.status) {
          case network.STATUS_NONE:
          case network.STATUS_PENDING:
            return null
          case network.STATUS_SUCCESS:
            return <Component/>
          case network.STATUS_FAILURE:
            return <Redirect to={{pathname: "/login", state: {from: location}}} />
        }
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
