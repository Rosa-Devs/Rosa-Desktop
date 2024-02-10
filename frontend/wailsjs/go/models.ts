export namespace manifest {
	
	export class Manifest {
	    name: string;
	    pubsub: string;
	    chiper: string;
	    optional: string;
	
	    static createFrom(source: any = {}) {
	        return new Manifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.pubsub = source["pubsub"];
	        this.chiper = source["chiper"];
	        this.optional = source["optional"];
	    }
	}

}

export namespace src {
	
	export class Message {
	    datatype: number;
	    sender: string;
	    data: string;
	    time: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.datatype = source["datatype"];
	        this.sender = source["sender"];
	        this.data = source["data"];
	        this.time = source["time"];
	    }
	}

}

