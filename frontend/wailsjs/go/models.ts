export namespace dto {
	
	export class ScannedPhoto {
	    filename: string;
	    path: string;
	    relativePath: string;
	    extension: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ScannedPhoto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.path = source["path"];
	        this.relativePath = source["relativePath"];
	        this.extension = source["extension"];
	        this.size = source["size"];
	    }
	}

}

export namespace model {
	
	export class Setting {
	    ID: number;
	    Source: string;
	    Dest: string;
	
	    static createFrom(source: any = {}) {
	        return new Setting(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Source = source["Source"];
	        this.Dest = source["Dest"];
	    }
	}

}

