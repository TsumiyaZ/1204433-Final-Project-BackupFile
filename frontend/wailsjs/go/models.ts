export namespace dto {
	
	export class GalleryPhoto {
	    id: number;
	    filename: string;
	    path: string;
	    size: number;
	    modifiedAt: string;
	    description: string;
	    status: string;
	    analysisStatus: string;
	    analysisError: string;
	
	    static createFrom(source: any = {}) {
	        return new GalleryPhoto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.filename = source["filename"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.modifiedAt = source["modifiedAt"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.analysisStatus = source["analysisStatus"];
	        this.analysisError = source["analysisError"];
	    }
	}
	export class MoveItemResult {
	    sourcePath: string;
	    filename: string;
	    destinationPath: string;
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new MoveItemResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourcePath = source["sourcePath"];
	        this.filename = source["filename"];
	        this.destinationPath = source["destinationPath"];
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}
	export class MovePhotosRequest {
	    paths: string[];
	
	    static createFrom(source: any = {}) {
	        return new MovePhotosRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.paths = source["paths"];
	    }
	}
	export class MovePhotosResult {
	    requested: number;
	    moved: number;
	    skipped: number;
	    failed: number;
	    durationMs: number;
	    items: MoveItemResult[];
	
	    static createFrom(source: any = {}) {
	        return new MovePhotosResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.requested = source["requested"];
	        this.moved = source["moved"];
	        this.skipped = source["skipped"];
	        this.failed = source["failed"];
	        this.durationMs = source["durationMs"];
	        this.items = this.convertValues(source["items"], MoveItemResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PhotoSearchResult {
	    id: number;
	    filename: string;
	    score: number;
	
	    static createFrom(source: any = {}) {
	        return new PhotoSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.filename = source["filename"];
	        this.score = source["score"];
	    }
	}
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

export namespace main {
	
	export class EmbeddingStatus {
	    state: string;
	    message: string;
	    progress: number;
	    progressVisible: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EmbeddingStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.state = source["state"];
	        this.message = source["message"];
	        this.progress = source["progress"];
	        this.progressVisible = source["progressVisible"];
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

