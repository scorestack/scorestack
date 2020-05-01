import React from 'react';
import {
  EuiPage,
  EuiPageBody,
  EuiPageContent,
  EuiPageContentHeader,
  EuiPageContentHeaderSection,
  EuiPageContentBody,
  EuiPageSideBar,
  EuiSideNav,
  EuiText,
  EuiTitle,
} from '@elastic/eui';
import { Attribute } from './attribute';

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
    };
  }

  componentDidMount() {
    const { httpClient } = this.props;
    httpClient.get('../api/scorestack/attribute').then((resp) => {
      this.setState({ checks: resp.data });
      const navItems = [];
      let itemId = 0;
      for (const group of Object.keys(this.state.checks)) {
        const subItems = [];
        for (const check of Object.keys(this.state.checks[group])) {
          subItems.push({
            name: this.state.checks[group][check].name,
            id: itemId,
            onClick: () => {
              this.setState({
                currentCheck: function () {
                  const attributes = Object.keys(this.props.attributes).map((key) => (
                    <Attribute
                      key={`${key}-${this.props.id}`}
                      id={this.props.id}
                      name={key}
                      value={this.props.attributes[key]}
                      client={this.props.httpClient}
                    />
                  ));
                  return (
                    <div>
                      <EuiPageContentHeader>
                        <EuiPageContentHeaderSection>
                          <EuiTitle>
                            <h2>{this.props.name}</h2>
                          </EuiTitle>
                        </EuiPageContentHeaderSection>
                      </EuiPageContentHeader>
                      <EuiPageContentBody>
                        {attributes}
                      </EuiPageContentBody>
                    </div>
                  );
                }
              });
            },
          });
          itemId++;
        }
        navItems.push({
          name: group,
          id: itemId,
          items: subItems,
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
            items={this.state.navItems}
          />
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
