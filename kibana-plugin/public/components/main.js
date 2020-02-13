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
import { Check } from './check';

export class Main extends React.Component {
  constructor(props) {
    super(props);

    this.state = { attributes: {} };
  }

  componentDidMount() {
    const { httpClient } = this.props;
    httpClient.get('../api/scorestack/attribute').then((resp) => {
      console.log(resp);
      this.setState({ attributes: resp.data["ssh-example"].attributes });
    });
  }

  render() {
    return (
      <EuiPage>
        <EuiPageBody>
          <EuiPageHeader>
            <EuiTitle size="l">
              <h1>Check Attributes</h1>
            </EuiTitle>
          </EuiPageHeader>
          <EuiPageContent>
            <Check attribs={this.state.attribs} id="ssh-example" name="Example SSH Check" httpClient={this.props.httpClient} />
          </EuiPageContent>
        </EuiPageBody>
      </EuiPage>
    );
  }
}