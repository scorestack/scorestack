import React from 'react';
import {
	EuiTitle,
	EuiPageContentBody,
	EuiPageContentHeader,
	EuiPageContentHeaderSection,
} from '@elastic/eui';
import { Attribute } from './attribute';

export class Check extends React.Component {
	constructor(props) {
		super(props);
	}

	render() {
		const attributes = Object.keys(this.props.attributes).map((key) =>
			<Attribute
				key={key}
				id={this.props.id}
				name={key}
				value={this.props.attributes[key]}
				client={this.props.httpClient} />
		);
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
}