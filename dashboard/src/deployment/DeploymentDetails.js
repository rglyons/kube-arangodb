import ReactTimeout from 'react-timeout';
import React, { Component } from 'react';
import api from '../api/api.js';
import Loading from '../util/Loading.js';
import MemberList from './MemberList.js';
import styled from 'react-emotion';
import { Loader } from 'semantic-ui-react';

const LoaderBox = styled('span')`
  float: right;
  width: 0;
  padding-right: 1em;
  margin-right: 1em;
  margin-top: 1em;
  max-width: 0;
  display: inline-block;
`;

const MemberGroupsView = ({memberGroups, namespace}) => (
  <div>
    {memberGroups.map((item) => <MemberList 
      key={item.group}
      group={item.group}
      members={item.members}
      namespace={namespace}
    />)}
  </div>
);

class DeploymentDetails extends Component {
  state = {
    loading: true,
    error: undefined
  };

  componentDidMount() {
    this.reloadDeployment();
  }

  reloadDeployment = async() => {
    try {
      this.setState({loading:true});
      const result = await api.get(`/api/deployment/${this.props.name}`);
      this.setState({
        deployment: result,
        loading: false,
        error: undefined
      });
    } catch (e) {
      this.setState({
        loading: false,
        error: e.message
      });
    }
    this.props.setTimeout(this.reloadDeployment, 5000);
  }

  render() {
    const d = this.state.deployment;
    if (!d) {
      return (<Loading/>);
    }
    return (
      <div>
        <LoaderBox><Loader size="mini" active={this.state.loading} inline/></LoaderBox>
        <MemberGroupsView memberGroups={d.member_groups} namespace={d.namespace}/>
      </div>
    );
  }
}

export default ReactTimeout(DeploymentDetails);