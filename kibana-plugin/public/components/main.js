import React from 'react';
import {
  EuiPage,
  EuiPageHeader,
  EuiTitle,
  EuiPageBody,
  EuiPageContent,
  EuiPageSideBar,
  EuiSideNav,
  EuiSideNavItem,
  EuiPageContentHeader,
  EuiPageContentBody,
  EuiFlexGroup,
  EuiFlexItem,
  EuiFormRow,
  EuiButton,
  EuiFieldText,
  EuiPopover,
  EuiForm,
  EuiFieldNumber,
  EuiRange,
  EuiSpacer,
  EuiSwitch,
  EuiText,
  EuiButtonIcon,
} from '@elastic/eui';
import { Check } from './check';

export class Main extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      checks: {},
      currentCheck: <EuiText>Click on any of your checks to the left to configure their attributes.</EuiText>,
      navItems: [{
        name: 'Loading...',
        id: 0,
      }],
      selectedItemName: 'Loading...',
    }
  }

  selectItem = name => {
    this.setState({
      selectedItemName: name,
    });
  };

  createItem = (name, data = {}) => {
    // NOTE: Duplicate `name` values will cause `id` collisions.
    return {
      ...data,
      id: name,
      name,
      isSelected: this.state.selectedItemName === name,
      onClick: () => this.selectItem(name),
    };
  };

  componentDidMount() {
    const { httpClient } = this.props;
    httpClient.get('../api/scorestack/attribute').then((resp) => {
      this.setState({ checks: resp.data });
      let navItems = []
      let itemId = 0;
      for (let group of Object.keys(this.state.checks)) {
        let subItems = []
        for (let check of Object.keys(this.state.checks[group])) {
          subItems.push({
            name: this.state.checks[group][check].name,
            id: itemId,
            onClick: () => {
              this.setState({
                curentCheck: <Check
                  id={check}
                  name={this.state.checks[group][check].name}
                  attributes={this.state.checks[group][check].attributes}
                  httpClient={this.props.httpClient} />
              });
            },
          })
          itemId++;
        }
        navItems.push(this.createItem(group, {
          items: subItems,
        }));
      }
      this.setState({ navItems: navItems });
    });
  }

  render() {
    return (
      <EuiPage>
        <EuiPageSideBar>
          <EuiSideNav
            mobileTitle="Checks"
            items={this.state.navItems} />
        </EuiPageSideBar>
        <EuiPageBody>
          <EuiPageContent>
            {this.state.currentCheck}
          </EuiPageContent>
        </EuiPageBody>
      </EuiPage >
    );
  }
}