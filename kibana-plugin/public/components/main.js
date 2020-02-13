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
      }]
    }
  }

  componentDidMount() {
    const { httpClient } = this.props;
    httpClient.get('../api/scorestack/attribute').then((resp) => {
      console.log(resp);
      this.setState({ checks: resp.data });
      let navItems = []
      let itemId = 0;
      for (let check of Object.keys(this.state.checks)) {
        navItems.push({
          name: this.state.checks[check].name,
          id: itemId,
          onClick: () => {
            this.setState({
              currentCheck: <Check
                id={check}
                name={this.state.checks[check].name}
                attributes={this.state.checks[check].attributes}
                httpClient={this.props.httpClient} />
            });
          },
        });
        itemId++;
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