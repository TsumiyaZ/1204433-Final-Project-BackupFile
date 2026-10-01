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

