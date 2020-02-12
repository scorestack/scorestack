import React from 'react';
import {
  EuiFlexGroup,
} from '@elastic/eui';

export class Attribute extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      isShown: false,
      value: props.value,
      formValue: '',
      isLoading: false,
    };
  }

  onSaveButtonClick = () => {
    this.setState({
      isLoading: true,
    });
    const { httpClient } = this.props.client;
    httpClient.post(`../api/scorestack/attribute/${this.props.id}/${this.props.key}`, JSON.stringify({
      'value': this.state.formValue,
    }, { headers: { 'Content-Type': 'application/json' } })).them((resp) => {
      this.setState({
        isLoading: false,
      });
      this.state.value = this.state.formValue;
      this.state.formValue = '';
    });
  }

  onShowButtonClick = () => {
    this.setState({
      isShown: !this.state.isShown,
    });
  };

  hideValue = () => {
    this.setState({
      isShown: false,
    });
  };

  showButton = (<EuiButtonIcon iconType='eye' onclick={this.onShowButtonClick} />);

  saveButton = (<EuiButton isLoading={this.state.isLoading} onclick={this.onSaveButtonClick}>Save</EuiButton>)

  render() {
    return (
      <EuiFlexGroup style={{ maxWidth: 600 }}>
        <EuiFlexItem>;
          <EuiFormRow label={this.props.key}>
            <EuiFieldText />
          </EuiFormRow>
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiPopover
              id="showValue"
              ownFocus
              button={this.showButton}
              isOpen={this.state.isShown}
              closePopover={this.hideValue.bind(this)}>
              <EuiText value={this.state.formValue}>
                <code>{this.state.value}</code>
              </EuiText>
            </EuiPopover>
          </EuiFormRow>
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            {this.saveButton}
          </EuiFormRow>
        </EuiFlexItem>
      </EuiFlexGroup>
    )
  }

}