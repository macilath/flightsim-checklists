import React, { ChangeEvent } from 'react';
import { Dialog, DialogTitle, Button, Input } from '@material-ui/core';

export interface NewAircraft {
    name: string
    alias: string
}

interface AircraftFormState {
    touched: boolean
    newName: string
    newAlias: string
}

interface ModalProps {
    open: boolean
    onClose: any
    onSubmit: any
}

export default class AddAircraftDialog extends React.Component<ModalProps, AircraftFormState> {
    constructor(props: ModalProps) {
        super(props);
        this.state = {
            touched: false,
            newName: '',
            newAlias: ''
        };
    }

    onNameChange(e: ChangeEvent<HTMLInputElement>) {
        this.setState({
            newName: e.target.value,
            touched: true
        });
    }

    onAliasChange(e: ChangeEvent<HTMLInputElement>) {
        this.setState({
            newAlias: e.target.value,
            touched: true
        });
    }

    submitAndClose() {
        const ac: NewAircraft = {
            name: this.state.newName,
            alias: this.state.newAlias
        };
        this.props.onSubmit(ac);
    }

    cancelOut() {
        this.setState({
            touched: false,
            newName: '',
            newAlias: ''
        }, () => this.props.onClose());
    }

    render() {
        return (
            <Dialog open={this.props.open} onClose={this.props.onClose}>
                <DialogTitle>Add an Aircraft</DialogTitle>
                <Input value={this.state.newName} onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.onNameChange(e)} />
                <Input value={this.state.newAlias} onChange={(e: React.ChangeEvent<HTMLInputElement>) => this.onAliasChange(e)} />
                <Button onClick={() => this.submitAndClose()}>Submit</Button><Button onClick={() => this.cancelOut()}>Cancel</Button>
            </Dialog>
        )
    }
}
