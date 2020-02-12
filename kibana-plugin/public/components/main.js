import React from 'react';
import {
  EuiPage,
  EuiPageHeader,
  EuiTitle,
  EuiPageBody,
  EuiPageContent,
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
import Attribute from './attribute';

export class Main extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      isAttributeShown: false,
      attribs = {},
    };
  }

  componentDidMount() {
    const { httpClient } = this.props;
    httpClient.get('../api/scorestack/attribute').then((resp) => {
      console.log(resp);
      this.setState({ attribs: resp.data["ssh-example"] });
    });
  }

  render() {
    const attributeItems = Object.keys(this.state.attribs).map((key) => {
      <Attribute id="ssh-example" key={key} value={this.state.attribs[key] || 'Loading...'} client={this.props.httpClient} />
    })
    return (
      <EuiPage>
        <EuiPageBody>
          <EuiPageHeader>
            <EuiTitle size="l">
              <h1>Check Attributes</h1>
            </EuiTitle>
          </EuiPageHeader>
          <EuiPageContent>
            {attributeItems}
          </EuiPageContent>
        </EuiPageBody>
      </EuiPage>
    );
  }
}