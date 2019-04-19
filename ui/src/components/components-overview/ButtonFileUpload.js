import React from "react";

class ButtonFileUpload extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            file: null,
            fileUrl: "",
            callBack: props.callBack,
        };

        this.getFile = this.getFile.bind(this);
    }

    pressButton(e) {
        e.preventDefault();
    }

    getFile(e) {
        e.preventDefault();
        console.log('Hello world');

        let reader = new FileReader();
        let file = e.target.files[0];

        reader.onloadend = () => {
            this.setState({
                file: file,
                fileUrl: reader.result
            });

            console.log('handle uploading-', this.state.file);
            this.state.callBack(this.state.file, this.state.fileUrl);
        }

        reader.readAsText(file);
    }

    render() {
        return (
            <div className="custom-file mb-3">
                <input type="file" className="custom-file-input" id="customFile2" onChange={this.getFile} />
                <label className="custom-file-label" htmlFor="customFile2">
                    Choose file...
                </label>
            </div>
        );
    }
}

export default ButtonFileUpload;

