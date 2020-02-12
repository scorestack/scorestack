import React from 'react';
import {
  EuiFlexGroup,
  EuiFlexItem,
  EuiButtonIcon,
  EuiButton,
  EuiFormRow,
  EuiFieldText,
  EuiPopover,
  EuiText,
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
    const httpClient = this.props.client;
    httpClient.post(`../api/scorestack/attribute/${this.props.id}/${this.props.name}`, JSON.stringify({
      'value': this.state.formValue,
    }, { headers: { 'Content-Type': 'application/json' } })).then((resp) => {
      console.log(resp);
      this.setState({
        isLoading: false,
        value: this.state.formValue,
      });
      this.setState({ formValue: '' });
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

  onChange = e => {
    this.setState({
      formValue: e.target.value,
    });
  };

  render() {
    const showButton = (<EuiButtonIcon iconType='eye' onClick={this.onShowButtonClick} />);
    const saveButton = (<EuiButton isLoading={this.state.isLoading} onClick={this.onSaveButtonClick}>Save</EuiButton>)
    return (
      <EuiFlexGroup style={{ maxWidth: 600 }}>
        <EuiFlexItem>
          <EuiFormRow label={this.props.name}>
            <EuiFieldText value={this.state.formValue} onChange={this.onChange} />
          </EuiFormRow>
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            <EuiPopover
              id="showValue"
              ownFocus
              button={showButton}
              isOpen={this.state.isShown}
              closePopover={this.hideValue.bind(this)}>
              <EuiText>
                <code>{this.state.value}</code>
              </EuiText>
            </EuiPopover>
          </EuiFormRow>
        </EuiFlexItem>
        <EuiFlexItem grow={false}>
          <EuiFormRow hasEmptyLabelSpace>
            {saveButton}
          </EuiFormRow>
        </EuiFlexItem>
      </EuiFlexGroup>
    )
  }

}