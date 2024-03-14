export namespace manifest {
	
	export class Manifest {
	    name: string;
	    uid: string;
	    pubsub: string;
	    chiper: string;
	    optional: string;
	
	    static createFrom(source: any = {}) {
	        return new Manifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.uid = source["uid"];
	        this.pubsub = source["pubsub"];
	        this.chiper = source["chiper"];
	        this.optional = source["optional"];
	    }
	}

}

