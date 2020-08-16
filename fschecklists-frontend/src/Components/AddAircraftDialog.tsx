import React from 'react';
import { Dialog, DialogTitle, Button, Input } from '@material-ui/core';

interface ModalProps {
    open: boolean
    onClose: any
    onSubmit: any
}

export default class AddAircraftDialog extends React.Component<ModalProps, {}> {
    
    render() {
        let apple = '';
        return (
            <Dialog open={this.props.open} onClose={this.props.onClose}>
                <DialogTitle>Add an Aircraft</DialogTitle>
                <Input value={apple}></Input>
                <Button onClick={this.props.onSubmit(apple)}>Submit</Button><Button onClick={this.props.onClose}>Cancel</Button>
            </Dialog>
        )
    }
}
